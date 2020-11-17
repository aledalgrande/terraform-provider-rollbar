package rollbar

/*
 * This file contains constants for resource names and strings to use in logging
 * and error output.
 */

// Resource names
const (
	resNameNotificationsEmail = "rollbar_notifications_email_integration"
	resNameNotificationsSlack = "rollbar_notifications_slack_integration"
)

// Strings for human-readable log messages
const (
	resNotificationsEmail = resNameNotificationsEmail + " resource"
	resNotificationsSlack = resNameNotificationsSlack + " resource"
)
