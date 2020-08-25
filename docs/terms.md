# List of terms

Below is a complete list of terms used in adam.

## Errors

| **Term**                                                     | **Default**                                                  | **Description**                                              |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `error.title`                                                | Error                                                        | The title of an error message.                               |
| `errors.internal.description.default`                        | Oh no! Something went wrong and I couldn't finish executing your command. I've informed my team and they'll get on fixing the bug asap. | The default description of an internal error embed.          |
| `errors.restriction.description.default`                     | 👮 You are not allowed to use this command.                   | The default description of an restriction error embed.       |
| `errors.insufficient_bot_permissions.description.single`     | It seems as if I don't have sufficient permissions to run this command. Please give me the "`{{missing_permission}}`" permission and try again. | The description of the the insufficient bot permissions error embed, if only a single permission is missing.<br /><br />**Keys:**<br />- `missing_permission` - the name of the missing permission |
| `errors.insufficient_bot_permissions.description.multi`      | It seems as if I don't have sufficient permissions to run this command. Please give me the following "+    "permissions and try again. | The description of the the insufficient bot permissions error embed, if multiple permissions are missing. |
| `errors.insufficient_bot_permissions.missing_permissions.name` | Missing Permission                                           | The name of the field containing a list of missing permissions. |
| `errors.argument_parsing_error.reason.name`                  | Reason                                                       | The name of the field containing the reason for the error.   |
| `errors.error_id`                                            | Error-ID: {{error_id}}                                       | The footer of an error embed.<br /><br />**Keys:**<br /> - `error_id` - the sentry event id of the error |
| `info.title`                                                 | Info                                                         | The title of an info message.                                |

## Permissions

| **Term**                            | **Default**                          |
| ----------------------------------- | ------------------------------------ |
| `permissions.create_instant_invite` | Create Invite                        |
| `permissions.kick_members`          | Kick Members                         |
| `permissions.ban_members`           | Ban Members                          |
| `permissions.administrator`         | Administrator                        |
| `permissions.manage_channels`       | Manage Channels                      |
| `permissions.manage_guild`          | Manage Server                        |
| `permissions.add_reactions`         | Add Reactions                        |
| `permissions.view_audit_log`        | View Audit Log                       |
| `permissions.priority_speaker`      | Priority Speaker                     |
| `permissions.stream`                | Video                                |
| `permissions.view_channel`          | View Channel                         |
| `permissions.send_messages`         | Send Messages                        |
| `permissions.send_tts_messages`     | Send TTS Messages                    |
| `permissions.manage_messages`       | Manage Messages                      |
| `permissions.embed_links`           | Embed Links                          |
| `permissions.attach_files`          | Attach Files                         |
| `permissions.read_message_history`  | Read Message History                 |
| `permissions.use_external_emojis`   | Mention everyone, here and All Roles |
| `permissions.connect`               | Connect                              |
| `permissions.speak`                 | Speak                                |
| `permissions.mute_members`          | Mute Members                         |
| `permissions.deafen_members`        | Deafen Members                       |
| `permissions.move_members`          | Move Members                         |
| `permissions.use_vad`               | Use Voice Activity                   |
| `permissions.change_nickname`       | Change Nickname                      |
| `permissions.manage_nicknames`      | Manage Nicknames                     |
| `permissions.manage_roles`          | Manage Roles                         |
| `permissions.manage_webhooks`       | Manage Webhooks                      |
| `permissions.manage_emojis`         | Manage Emojis                        |