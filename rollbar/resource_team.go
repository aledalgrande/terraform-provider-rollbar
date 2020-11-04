package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
)

// resourceTeam constructs a resource representing a Rollbar team.
func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		DeleteContext: resourceTeamDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "standard",
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	level := d.Get("access_level").(string)
	l := log.With().Str("name", name).Str("access_level", level).Logger()
	l.Info().Msg("Creating rollbar_team resource")
	c := m.(*client.RollbarApiClient)
	t, err := c.CreateTeam(name, level)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	teamID := t.ID
	l = l.With().Int("teamID", teamID).Logger()
	d.SetId(strconv.Itoa(teamID))
	l.Debug().Int("id", teamID).Msg("Successfully created rollbar_team resource")
	return resourceTeamRead(ctx, d, m)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Reading Rollbar team from API")
	c := m.(*client.RollbarApiClient)
	t, err := c.ReadTeam(id)
	if err == client.ErrNotFound {
		return handleErrNotFound(d, "team")
	}
	if err != nil {
		l.Err(err).Msg("error reading team resource")
		return diag.FromErr(err)
	}
	var errs []error
	errs = append(errs, d.Set("name", t.Name))
	errs = append(errs, d.Set("account_id", t.AccountID))
	errs = append(errs, d.Set("access_level", t.AccessLevel))
	for _, err = range errs {
		if err != nil {
			l.Error().Interface("errs", errs).Send()
			return diag.FromErr(errs[0])
		}
	}
	l.Debug().Msg("Successfully read rollbar_team resource")
	return nil
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	l := log.With().Int("id", id).Logger()
	l.Info().Msg("Deleting team resource")
	c := m.(*client.RollbarApiClient)
	err = c.DeleteTeam(id)
	if err != nil {
		l.Err(err).Msg("Error deleting rollbar_team resource")
		return diag.FromErr(err)
	}
	l.Debug().Msg("Successfully deleted rollbar_team resource")
	return nil
}
