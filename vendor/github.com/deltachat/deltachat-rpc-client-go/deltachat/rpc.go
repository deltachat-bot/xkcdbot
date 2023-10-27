package deltachat

import (
	"context"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat/option"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat/transport"
)

// Delta Chat RPC client. This is the root of the API.
type Rpc struct {
	// Context to be used on calls to Transport.CallResult() and Transport.Call()
	Context   context.Context
	Transport transport.RpcTransport
}

// ---------------------------------------------
//  Misc top level functions
// ---------------------------------------------

// Check if an email address is valid.
func (self *Rpc) CheckEmailValidity(email string) (bool, error) {
	var valid bool
	err := self.Transport.CallResult(self.Context, &valid, "check_email_validity", email)
	return valid, err
}

// Get general system info.
func (self *Rpc) GetSystemInfo() (map[string]string, error) {
	var info map[string]string
	err := self.Transport.CallResult(self.Context, &info, "get_system_info")
	return info, err
}

// Get the next event.
func (self *Rpc) GetNextEvent() (AccountId, Event, error) {
	var event _Event
	err := self.Transport.CallResult(self.Context, &event, "get_next_event")
	if err != nil {
		return 0, nil, err
	}
	return event.ContextId, event.Event.ToEvent(), nil
}

// ---------------------------------------------
// Account Management
// ---------------------------------------------

// Create a new account.
func (self *Rpc) AddAccount() (AccountId, error) {
	var id AccountId
	err := self.Transport.CallResult(self.Context, &id, "add_account")
	return id, err
}

// Remove an account.
func (self *Rpc) RemoveAccount(accountId AccountId) error {
	return self.Transport.Call(self.Context, "remove_account", accountId)
}

// Return all available accounts.
func (self *Rpc) GetAllAccountIds() ([]AccountId, error) {
	var ids []AccountId
	err := self.Transport.CallResult(self.Context, &ids, "get_all_account_ids")
	return ids, err
}

// Select account id for internally selected state.
func (self *Rpc) SelectAccount(accountId AccountId) error {
	return self.Transport.Call(self.Context, "select_account", accountId)
}

// Get the selected account id of the internal state.
func (self *Rpc) GetSelectedAccountId() (option.Option[AccountId], error) {
	var id option.Option[AccountId]
	err := self.Transport.CallResult(self.Context, &id, "get_selected_account_id")
	return id, err
}

// TODO: get_all_accounts

// Start the I/O of all accounts.
func (self *Rpc) StartIoForAllAccounts() error {
	return self.Transport.Call(self.Context, "start_io_for_all_accounts")
}

// Stop the I/O of all accounts.
func (self *Rpc) StopIoForAllAccounts() error {
	return self.Transport.Call(self.Context, "stop_io_for_all_accounts")
}

// ---------------------------------------------
// Methods that work on individual accounts
// ---------------------------------------------

// Start the account I/O.
func (self *Rpc) StartIo(accountId AccountId) error {
	return self.Transport.Call(self.Context, "start_io", accountId)
}

// Stop the account I/O.
func (self *Rpc) StopIo(accountId AccountId) error {
	return self.Transport.Call(self.Context, "stop_io", accountId)
}

// TODO: get_account_info

// Get the combined filesize of an account in bytes.
func (self *Rpc) GetAccountFileSize(accountId AccountId) (uint64, error) {
	var size uint64
	err := self.Transport.CallResult(self.Context, &size, "get_account_file_size", accountId)
	return size, err
}

// TODO: get_provider_info

// Checks if the account is already configured.
func (self *Rpc) IsConfigured(accountId AccountId) (bool, error) {
	var configured bool
	err := self.Transport.CallResult(self.Context, &configured, "is_configured", accountId)
	return configured, err
}

// Get system info for an account.
func (self *Rpc) GetInfo(accountId AccountId) (map[string]string, error) {
	var info map[string]string
	err := self.Transport.CallResult(self.Context, &info, "get_info", accountId)
	return info, err
}

// Set account configuration value.
func (self *Rpc) SetConfig(accountId AccountId, key string, value option.Option[string]) error {
	return self.Transport.Call(self.Context, "set_config", accountId, key, value)
}

// Tweak several account configuration values in a batch.
func (self *Rpc) BatchSetConfig(accountId AccountId, config map[string]option.Option[string]) error {
	return self.Transport.Call(self.Context, "batch_set_config", accountId, config)
}

// TODO: set_config_from_qr
// TODO: check_qr

// Get custom UI-specific configuration value set with SetUiConfig().
func (self *Rpc) GetConfig(accountId AccountId, key string) (option.Option[string], error) {
	var value option.Option[string]
	err := self.Transport.CallResult(self.Context, &value, "get_config", accountId, key)
	return value, err
}

// Get a batch of account configuration values.
func (self *Rpc) BatchGetConfig(accountId AccountId, keys []string) (map[string]option.Option[string], error) {
	var values map[string]option.Option[string]
	err := self.Transport.CallResult(self.Context, &values, "batch_get_config", accountId, keys)
	return values, err
}

// Set stock strings.
func (self *Rpc) SetStockStrings(translations map[uint]string) error {
	return self.Transport.Call(self.Context, "set_stock_strings", translations)
}

// Configures an account with the currently set parameters.
// Setup the credential config before calling this.
func (self *Rpc) Configure(accountId AccountId) error {
	return self.Transport.Call(self.Context, "configure", accountId)
}

// Signal an ongoing process to stop.
func (self *Rpc) StopOngoingProcess(accountId AccountId) error {
	return self.Transport.Call(self.Context, "stop_ongoing_process", accountId)
}

// Export public and private keys to the specified directory.
// Note that the account does not have to be started.
func (self *Rpc) ExportSelfKeys(accountId AccountId, path string) error {
	return self.Transport.Call(self.Context, "export_self_keys", accountId, path, nil)
}

// Import private keys found in the specified directory.
func (self *Rpc) ImportSelfKeys(accountId AccountId, path string) error {
	return self.Transport.Call(self.Context, "import_self_keys", accountId, path, nil)
}

// Returns the message IDs of all fresh messages of any chat.
// Typically used for implementing notification summaries
// or badge counters e.g. on the app icon.
// The list is already sorted and starts with the most recent fresh message.
//
// Messages belonging to muted chats or to the contact requests are not returned;
// these messages should not be notified
// and also badge counters should not include these messages.
//
// To get the number of fresh messages for a single chat, muted or not,
// use GetFreshMsgCnt().
func (self *Rpc) GetFreshMsgs(accountId AccountId) ([]MsgId, error) {
	var ids []MsgId
	err := self.Transport.CallResult(self.Context, &ids, "get_fresh_msgs", accountId)
	return ids, err
}

// Get the number of fresh messages in a chat.
// Typically used to implement a badge with a number in the chatlist.
//
// If the specified chat is muted,
// the UI should show the badge counter "less obtrusive",
// e.g. using "gray" instead of "red" color.
func (self *Rpc) GetFreshMsgCnt(accountId AccountId, chatId ChatId) (uint, error) {
	var count uint
	err := self.Transport.CallResult(self.Context, &count, "get_fresh_msg_cnt", accountId, chatId)
	return count, err
}

// Gets messages to be processed by the bot and returns their IDs.
//
// Only messages with database ID higher than last_msg_id config value
// are returned. After processing the messages, the bot should
// update last_msg_id by calling MarkseenMsgs()
// or manually updating the value to avoid getting already
// processed messages.
func (self *Rpc) GetNextMsgs(accountId AccountId) ([]MsgId, error) {
	var ids []MsgId
	err := self.Transport.CallResult(self.Context, &ids, "get_next_msgs", accountId)
	return ids, err
}

// Waits for messages to be processed by the bot and returns their IDs.
//
// This function is similar to GetNextMsgs(),
// but waits for internal new message notification before returning.
// New message notification is sent when new message is added to the database,
// on initialization, when I/O is started and when I/O is stopped.
// This allows bots to use WaitNextMsgs() in a loop to process
// old messages after initialization and during the bot runtime.
// To shutdown the bot, stopping I/O can be used to interrupt
// pending or next WaitNextMsgs() call.
func (self *Rpc) WaitNextMsgs(accountId AccountId) ([]MsgId, error) {
	var ids []MsgId
	err := self.Transport.CallResult(self.Context, &ids, "wait_next_msgs", accountId)
	return ids, err
}

// Estimate the number of messages that will be deleted
// by the SetConfig()-options `delete_device_after` or `delete_server_after`.
// This is typically used to show the estimated impact to the user
// before actually enabling deletion of old messages.
func (self *Rpc) EstimateAutoDeletionCount(accountId AccountId, fromServer bool, seconds int64) (uint, error) {
	var count uint
	err := self.Transport.CallResult(self.Context, &count, "estimate_auto_deletion_count", accountId, fromServer, seconds)
	return count, err
}

// ---------------------------------------------
//  autocrypt
// ---------------------------------------------

// Start the AutoCrypt key transfer process.
func (self *Rpc) InitiateAutocryptKeyTransfer(accountId AccountId) (string, error) {
	var result string
	err := self.Transport.CallResult(self.Context, &result, "initiate_autocrypt_key_transfer", accountId)
	return result, err
}

// Continue the AutoCrypt key transfer process.
func (self *Rpc) ContinueAutocryptKeyTransfer(accountId AccountId, msgId MsgId, setupCode string) error {
	return self.Transport.Call(self.Context, "continue_autocrypt_key_transfer", accountId, msgId, setupCode)
}

// ---------------------------------------------
//   chat list
// ---------------------------------------------

func (self *Rpc) GetChatlistEntries(accountId AccountId, listFlags option.Option[uint], query option.Option[string], contactId option.Option[ContactId]) ([]ChatId, error) {
	var entries []ChatId
	err := self.Transport.CallResult(self.Context, &entries, "get_chatlist_entries", accountId, listFlags, query, contactId)
	return entries, err
}

func (self *Rpc) GetChatlistItemsByEntries(accountId AccountId, entries []ChatId) (map[ChatId]*ChatListItem, error) {
	var itemsMap map[ChatId]*ChatListItem
	err := self.Transport.CallResult(self.Context, &itemsMap, "get_chatlist_items_by_entries", accountId, entries)
	return itemsMap, err
}

// ---------------------------------------------
//  chat
// ---------------------------------------------

func (self *Rpc) GetFullChatById(accountId AccountId, chatId ChatId) (*FullChatSnapshot, error) {
	var result FullChatSnapshot
	err := self.Transport.CallResult(self.Context, &result, "get_full_chat_by_id", accountId, chatId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// get basic info about a chat,
// use GetFullChatById() instead if you need more information
func (self *Rpc) GetBasicChatInfo(accountId AccountId, chatId ChatId) (*BasicChatSnapshot, error) {
	var result BasicChatSnapshot
	err := self.Transport.CallResult(self.Context, &result, "get_basic_chat_info", accountId, chatId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (self *Rpc) AcceptChat(accountId AccountId, chatId ChatId) error {
	return self.Transport.Call(self.Context, "accept_chat", accountId, chatId)
}

func (self *Rpc) BlockChat(accountId AccountId, chatId ChatId) error {
	return self.Transport.Call(self.Context, "block_chat", accountId, chatId)
}

// Delete a chat.
//
// Messages are deleted from the device and the chat database entry is deleted.
// After that, the event #DC_EVENT_MSGS_CHANGED is posted.
//
// Things that are _not done_ implicitly:
//
//   - Messages are **not deleted from the server**.
//   - The chat or the contact is **not blocked**, so new messages from the user/the group may appear as a contact request
//     and the user may create the chat again.
//   - **Groups are not left** - this would
//     be unexpected as (1) deleting a normal chat also does not prevent new mails
//     from arriving, (2) leaving a group requires sending a message to
//     all group members - especially for groups not used for a longer time, this is
//     really unexpected when deletion results in contacting all members again,
//     (3) only leaving groups is also a valid usecase.
//
// To leave a chat explicitly, use leave_group()
func (self *Rpc) DeleteChat(accountId AccountId, chatId ChatId) error {
	return self.Transport.Call(self.Context, "delete_chat", accountId, chatId)
}

// Get encryption info for this chat.
// Get a multi-line encryption info, containing encryption preferences of all members.
// Can be used to find out why messages sent to group are not encrypted.
//
// returns Multi-line text
func (self *Rpc) GetChatEncryptionInfo(accountId AccountId, chatId ChatId) (string, error) {
	var data string
	err := self.Transport.CallResult(self.Context, &data, "get_chat_encryption_info", accountId, chatId)
	return data, err
}

// Get Join-Group QR code text and SVG data.
func (self *Rpc) GetChatSecurejoinQrCodeSvg(accountId AccountId, chatId option.Option[ChatId]) (string, string, error) {
	var data [2]string
	err := self.Transport.CallResult(self.Context, &data, "get_chat_securejoin_qr_code_svg", accountId, chatId)
	return data[0], data[1], err
}

// Continue a Setup-Contact or Verified-Group-Invite protocol started on another device.
func (self *Rpc) SecureJoin(accountId AccountId, qrdata string) (ChatId, error) {
	var id ChatId
	err := self.Transport.CallResult(self.Context, &id, "secure_join", accountId, qrdata)
	return id, err
}

func (self *Rpc) LeaveGroup(accountId AccountId, chatId ChatId) error {
	return self.Transport.Call(self.Context, "leave_group", accountId, chatId)
}

// Remove a member from a group.
func (self *Rpc) RemoveContactFromChat(accountId AccountId, chatId ChatId, contactId ContactId) error {
	return self.Transport.Call(self.Context, "remove_contact_from_chat", accountId, chatId, contactId)
}

// Add a member to a group.
func (self *Rpc) AddContactToChat(accountId AccountId, chatId ChatId, contactId ContactId) error {
	return self.Transport.Call(self.Context, "add_contact_to_chat", accountId, chatId, contactId)
}

// Get the contact IDs belonging to a chat.
//
//   - for normal chats, the function always returns exactly one contact,
//     DC_CONTACT_ID_SELF is returned only for SELF-chats.
//
//   - for group chats all members are returned, DC_CONTACT_ID_SELF is returned
//     explicitly as it may happen that oneself gets removed from a still existing
//     group
//
// - for broadcasts, all recipients are returned, DC_CONTACT_ID_SELF is not included
//
//   - for mailing lists, the behavior is not documented currently, we will decide on that later.
//     for now, the UI should not show the list for mailing lists.
//     (we do not know all members and there is not always a global mailing list address,
//     so we could return only SELF or the known members; this is not decided yet)
func (self *Rpc) GetChatContacts(accountId AccountId, chatId ChatId) ([]ContactId, error) {
	var ids []ContactId
	err := self.Transport.CallResult(self.Context, &ids, "get_chat_contacts", accountId, chatId)
	return ids, err
}

// Create a new group chat.
// After creation, the group has only self-contact as member and is in unpromoted state.
func (self *Rpc) CreateGroupChat(accountId AccountId, name string, protected bool) (ChatId, error) {
	var id ChatId
	err := self.Transport.CallResult(self.Context, &id, "create_group_chat", accountId, name, protected)
	return id, err
}

// Create a new broadcast list.
func (self *Rpc) CreateBroadcastList(accountId AccountId) (ChatId, error) {
	var id ChatId
	err := self.Transport.CallResult(self.Context, &id, "create_broadcast_list", accountId)
	return id, err
}

// Set group name.
func (self *Rpc) SetChatName(accountId AccountId, chatId ChatId, name string) error {
	return self.Transport.Call(self.Context, "set_chat_name", accountId, chatId, name)
}

// Set group profile image.
//
// If the group is already _promoted_ (any message was sent to the group),
// all group members are informed by a special status message that is sent automatically by this function.
//
// Sends out #DC_EVENT_CHAT_MODIFIED and #DC_EVENT_MSGS_CHANGED if a status message was sent.
//
// To find out the profile image of a chat, use dc_chat_get_profile_image()
//
// @param image_path Full path of the image to use as the group image. The image will immediately be copied to the
//
//	`blobdir`; the original image will not be needed anymore.
//	 If you pass null here, the group image is deleted (for promoted groups, all members are informed about
//	 this change anyway).
func (self *Rpc) SetChatProfileImage(accountId AccountId, chatId ChatId, path option.Option[string]) error {
	return self.Transport.Call(self.Context, "set_chat_profile_image", accountId, chatId, path)
}

func (self *Rpc) SetChatVisibility(accountId AccountId, chatId ChatId, visibility ChatVisibility) error {
	return self.Transport.Call(self.Context, "set_chat_visibility", accountId, chatId, visibility)
}

func (self *Rpc) SetChatEphemeralTimer(accountId AccountId, chatId ChatId, timer uint) error {
	return self.Transport.Call(self.Context, "set_chat_ephemeral_timer", accountId, chatId, timer)
}

func (self *Rpc) GetChatEphemeralTimer(accountId AccountId, chatId ChatId) (uint, error) {
	var timer uint
	err := self.Transport.CallResult(self.Context, &timer, "get_chat_ephemeral_timer", accountId, chatId)
	return timer, err
}

// for now only text messages, because we only used text messages in desktop thusfar
func (self *Rpc) AddDeviceMessage(accountId AccountId, label string, text string) (MsgId, error) {
	var id MsgId
	err := self.Transport.CallResult(self.Context, &id, "add_device_message", accountId, label, text)
	return id, err
}

// Mark all messages in a chat as _noticed_.
// _Noticed_ messages are no longer _fresh_ and do not count as being unseen
// but are still waiting for being marked as "seen" using markseen_msgs()
// (IMAP/MDNs is not done for noticed messages).
//
// Calling this function usually results in the event #DC_EVENT_MSGS_NOTICED.
// See also markseen_msgs().
func (self *Rpc) MarknoticedChat(accountId AccountId, chatId ChatId) error {
	return self.Transport.Call(self.Context, "marknoticed_chat", accountId, chatId)
}

func (self *Rpc) GetFirstUnreadMessageOfChat(accountId AccountId, chatId ChatId) (option.Option[MsgId], error) {
	var id option.Option[MsgId]
	err := self.Transport.CallResult(self.Context, &id, "get_first_unread_message_of_chat", accountId, chatId)
	return id, err
}

// TODO: set_chat_mute_duration
// TODO: is_chat_muted

// ---------------------------------------------
// message list
// ---------------------------------------------

// Mark messages as presented to the user.
// Typically, UIs call this function on scrolling through the message list,
// when the messages are presented at least for a little moment.
// The concrete action depends on the type of the chat and on the users settings
// (dc_msgs_presented() may be a better name therefore, but well. :)
//
//   - For normal chats, the IMAP state is updated, MDN is sent
//     (if set_config()-options `mdns_enabled` is set)
//     and the internal state is changed to @ref DC_STATE_IN_SEEN to reflect these actions.
//
//   - For contact requests, no IMAP or MDNs is done
//     and the internal state is not changed therefore.
//     See also marknoticed_chat().
//
// Moreover, timer is started for incoming ephemeral messages.
// This also happens for contact requests chats.
//
// This function updates `last_msg_id` configuration value
// to the maximum of the current value and IDs passed to this function.
// Bots which mark messages as seen can rely on this side effect
// to avoid updating `last_msg_id` value manually.
//
// One #DC_EVENT_MSGS_NOTICED event is emitted per modified chat.
func (self *Rpc) MarkseenMsgs(accountId AccountId, msgIds []MsgId) error {
	return self.Transport.Call(self.Context, "markseen_msgs", accountId, msgIds)
}

func (self *Rpc) GetMessageIds(accountId AccountId, chatId ChatId, infoOnly, addDaymarker bool) ([]MsgId, error) {
	var ids []MsgId
	err := self.Transport.CallResult(self.Context, &ids, "get_message_ids", accountId, chatId, infoOnly, addDaymarker)
	return ids, err
}

// TODO: get_message_list_items

// Return map of this account configuration parameters.
func (self *Rpc) GetMessage(accountId AccountId, msgId MsgId) (*MsgSnapshot, error) {
	var snapshot MsgSnapshot
	err := self.Transport.CallResult(self.Context, &snapshot, "get_message", accountId, msgId)
	return &snapshot, err
}

// Get the HTML part of this message.
func (self *Rpc) GetMessageHtml(accountId AccountId, msgId MsgId) (option.Option[string], error) {
	var html option.Option[string]
	err := self.Transport.CallResult(self.Context, &html, "get_message_html", accountId, msgId)
	return html, err
}

// TODO: get_messages
// TODO: get_message_notification_info

// Delete messages. The messages are deleted on the current device and
// on the IMAP server.
func (self *Rpc) DeleteMessages(accountId AccountId, msgIds []MsgId) error {
	return self.Transport.Call(self.Context, "delete_messages", accountId, msgIds)
}

// Get an informational text for a single message. The text is multiline and may
// contain e.g. the raw text of the message.
//
// The max. text returned is typically longer (about 100000 characters) than the
// max. text returned by dc_msg_get_text() (about 30000 characters).
func (self *Rpc) GetMessageInfo(accountId AccountId, msgId MsgId) (string, error) {
	var info string
	err := self.Transport.CallResult(self.Context, &info, "get_message_info", accountId, msgId)
	return info, err
}

// Asks the core to start downloading a message fully.
// This function is typically called when the user hits the "Download" button
// that is shown by the UI in case `download_state` is `'Available'` or `'Failure'`
//
// On success, the @ref DC_MSG "view type of the message" may change
// or the message may be replaced completely by one or more messages with other message IDs.
// That may happen e.g. in cases where the message was encrypted
// and the type could not be determined without fully downloading.
// Downloaded content can be accessed as usual after download.
//
// To reflect these changes a @ref DC_EVENT_MSGS_CHANGED event will be emitted.
func (self *Rpc) DownloadFullMessage(accountId AccountId, msgId MsgId) error {
	return self.Transport.Call(self.Context, "download_full_message", accountId, msgId)
}

// Search messages containing the given query string.
// Searching can be done globally (chat_id=None) or in a specified chat only (chat_id set).
//
// Global search results are typically displayed using dc_msg_get_summary(), chat
// search results may just highlight the corresponding messages and present a
// prev/next button.
//
// For the global search, the result is limited to 1000 messages,
// this allows an incremental search done fast.
// So, when getting exactly 1000 messages, the result actually may be truncated;
// the UIs may display sth. like "1000+ messages found" in this case.
// The chat search (if chat_id is set) is not limited.
func (self *Rpc) SearchMessages(accountId AccountId, query string, chatId option.Option[ChatId]) ([]MsgId, error) {
	var msgIds []MsgId
	err := self.Transport.CallResult(self.Context, &msgIds, "search_messages", accountId, query, chatId)
	return msgIds, err
}

func (self *Rpc) MessageIdsToSearchResults(accountId AccountId, msgIds []MsgId) (map[MsgId]*MsgSearchResult, error) {
	var results map[MsgId]*MsgSearchResult
	err := self.Transport.CallResult(self.Context, &results, "message_ids_to_search_results", accountId, msgIds)
	return results, err
}

// ---------------------------------------------
//  contact
// ---------------------------------------------

// Get the properties of a single contact by ID.
func (self *Rpc) GetContact(accountId AccountId, contactId ContactId) (*ContactSnapshot, error) {
	var snapshot ContactSnapshot
	err := self.Transport.CallResult(self.Context, &snapshot, "get_contact", accountId, contactId)
	return &snapshot, err
}

// Add a single contact as a result of an explicit user action.
//
// Returns contact id of the created or existing contact
func (self *Rpc) CreateContact(accountId AccountId, email string, name string) (ContactId, error) {
	var id ContactId
	err := self.Transport.CallResult(self.Context, &id, "create_contact", accountId, email, name)
	return id, err
}

// Returns contact id of the created or existing DM chat with that contact
func (self *Rpc) CreateChatByContactId(accountId AccountId, contactId ContactId) (ChatId, error) {
	var id ChatId
	err := self.Transport.CallResult(self.Context, &id, "create_chat_by_contact_id", accountId, contactId)
	return id, err
}

func (self *Rpc) BlockContact(accountId AccountId, contactId ContactId) error {
	return self.Transport.Call(self.Context, "block_contact", accountId, contactId)
}

func (self *Rpc) UnblockContact(accountId AccountId, contactId ContactId) error {
	return self.Transport.Call(self.Context, "unblock_contact", accountId, contactId)
}

func (self *Rpc) GetBlockedContacts(accountId AccountId) ([]*ContactSnapshot, error) {
	var contacts []*ContactSnapshot
	err := self.Transport.CallResult(self.Context, &contacts, "get_blocked_contacts", accountId)
	return contacts, err
}

func (self *Rpc) GetContactIds(accountId AccountId, listFlags uint, query option.Option[string]) ([]ContactId, error) {
	var ids []ContactId
	err := self.Transport.CallResult(self.Context, &ids, "get_contact_ids", accountId, listFlags, query)
	return ids, err
}

// TODO: get_contacts
// TODO: get_contacts_by_ids

func (self *Rpc) DeleteContact(accountId AccountId, contactId ContactId) error {
	return self.Transport.Call(self.Context, "delete_contact", accountId, contactId)
}

func (self *Rpc) ChangeContactName(accountId AccountId, contactId ContactId, name string) error {
	return self.Transport.Call(self.Context, "change_contact_name", accountId, contactId, name)
}

// Get encryption info for a contact.
// Get a multi-line encryption info, containing your fingerprint and the
// fingerprint of the contact, used e.g. to compare the fingerprints for a simple out-of-band verification.
func (self *Rpc) GetContactEncryptionInfo(accountId AccountId, contactId ContactId) (string, error) {
	var data string
	err := self.Transport.CallResult(self.Context, &data, "get_contact_encryption_info", accountId, contactId)
	return data, err
}

// Check if an e-mail address belongs to a known and unblocked contact.
// To get a list of all known and unblocked contacts, use contacts_get_contacts().
//
// To validate an e-mail address independently of the contact database
// use check_email_validity().
func (self *Rpc) LookupContactIdByAddr(accountId AccountId, addr string) (option.Option[ContactId], error) {
	var id option.Option[ContactId]
	err := self.Transport.CallResult(self.Context, &id, "lookup_contact_id_by_addr", accountId, addr)
	return id, err
}

// ---------------------------------------------
//                   chat
// ---------------------------------------------

// Returns all message IDs of the given types in a chat.
// Typically used to show a gallery.
//
// The list is already sorted and starts with the oldest message.
// Clients should not try to re-sort the list as this would be an expensive action
// and would result in inconsistencies between clients.
//
// Setting `chat_id` to `None` (`null` in typescript) means get messages with media
// from any chat of the currently used account.
func (self *Rpc) GetChatMedia(accountId AccountId, chatId ChatId, messageType MsgType, orMessageType2 option.Option[MsgType], orMessageType3 option.Option[MsgType]) ([]MsgId, error) {
	var ids []MsgId
	err := self.Transport.CallResult(self.Context, &ids, "get_chat_media", accountId, chatId, messageType, orMessageType2, orMessageType3)
	return ids, err
}

// TODO: get_neighboring_chat_media

// ---------------------------------------------
//                   backup
// ---------------------------------------------

// Export account backup.
func (self *Rpc) ExportBackup(accountId AccountId, destination string, passphrase option.Option[string]) error {
	return self.Transport.Call(self.Context, "export_backup", accountId, destination, passphrase)
}

// Import account backup.
func (self *Rpc) ImportBackup(accountId AccountId, path string, passphrase option.Option[string]) error {
	return self.Transport.Call(self.Context, "import_backup", accountId, path, passphrase)
}

// Offers a backup for remote devices to retrieve.
//
// Can be cancelled by stopping the ongoing process.  Success or failure can be tracked
// via the `ImexProgress` event which should either reach `1000` for success or `0` for
// failure.
//
// This **stops IO** while it is running.
//
// Returns once a remote device has retrieved the backup, or is cancelled.
func (self *Rpc) ProvideBackup(accountId AccountId) error {
	return self.Transport.Call(self.Context, "provide_backup", accountId)
}

// Returns the text of the QR code for the running [`CommandApi::provide_backup`].
//
// This QR code text can be used in [`CommandApi::get_backup`] on a second device to
// retrieve the backup and setup this second device.
//
// This call will fail if there is currently no concurrent call to
// [`CommandApi::provide_backup`].  This call may block if the QR code is not yet
// ready.
func (self *Rpc) GetBackupQr(accountId AccountId) (string, error) {
	var result string
	err := self.Transport.CallResult(self.Context, &result, "get_backup_qr", accountId)
	return result, err
}

// Returns the rendered QR code for the running [`CommandApi::provide_backup`].
//
// This QR code can be used in [`CommandApi::get_backup`] on a second device to
// retrieve the backup and setup this second device.
//
// This call will fail if there is currently no concurrent call to
// [`CommandApi::provide_backup`].  This call may block if the QR code is not yet
// ready.
//
// Returns the QR code rendered as an SVG image.
func (self *Rpc) GetBackupQrSvg(accountId AccountId) (string, error) {
	var result string
	err := self.Transport.CallResult(self.Context, &result, "get_backup_qr_svg", accountId)
	return result, err
}

// Gets a backup from a remote provider.
//
// This retrieves the backup from a remote device over the network and imports it into
// the current device.
//
// Can be cancelled by stopping the ongoing process.
func (self *Rpc) GetBackup(accountId AccountId, qrText string) error {
	return self.Transport.Call(self.Context, "get_backup", accountId, qrText)
}

// ---------------------------------------------
//                connectivity
// ---------------------------------------------

// Indicate that the network likely has come back.
// or just that the network conditions might have changed
func (self *Rpc) MaybeNetwork() error {
	return self.Transport.Call(self.Context, "maybe_network")
}

// Get the current connectivity, i.e. whether the device is connected to the IMAP server.
// One of:
// - DC_CONNECTIVITY_NOT_CONNECTED (1000-1999): Show e.g. the string "Not connected" or a red dot
// - DC_CONNECTIVITY_CONNECTING (2000-2999): Show e.g. the string "Connecting…" or a yellow dot
// - DC_CONNECTIVITY_WORKING (3000-3999): Show e.g. the string "Getting new messages" or a spinning wheel
// - DC_CONNECTIVITY_CONNECTED (>=4000): Show e.g. the string "Connected" or a green dot
//
// We don't use exact values but ranges here so that we can split up
// states into multiple states in the future.
//
// Meant as a rough overview that can be shown
// e.g. in the title of the main screen.
//
// If the connectivity changes, a #DC_EVENT_CONNECTIVITY_CHANGED will be emitted.
func (self *Rpc) GetConnectivity(accountId AccountId) (uint, error) {
	var info uint
	err := self.Transport.CallResult(self.Context, &info, "get_connectivity", accountId)
	return info, err
}

// Get an overview of the current connectivity, and possibly more statistics.
// Meant to give the user more insight about the current status than
// the basic connectivity info returned by get_connectivity(); show this
// e.g., if the user taps on said basic connectivity info.
//
// If this page changes, a #DC_EVENT_CONNECTIVITY_CHANGED will be emitted.
//
// This comes as an HTML from the core so that we can easily improve it
// and the improvement instantly reaches all UIs.
func (self *Rpc) GetConnectivityHtml(accountId AccountId) (string, error) {
	var html string
	err := self.Transport.CallResult(self.Context, &html, "get_connectivity_html", accountId)
	return html, err
}

// ---------------------------------------------
//                  locations
// ---------------------------------------------

// TODO: get_locations

// ---------------------------------------------
//                   webxdc
// ---------------------------------------------

func (self *Rpc) SendWebxdcStatusUpdate(accountId AccountId, msgId MsgId, update string, description string) error {
	return self.Transport.Call(self.Context, "send_webxdc_status_update", accountId, msgId, update, description)
}

func (self *Rpc) GetWebxdcStatusUpdates(accountId AccountId, msgId MsgId, lastKnownSerial uint) (string, error) {
	var data string
	err := self.Transport.CallResult(self.Context, &data, "get_webxdc_status_updates", accountId, msgId, lastKnownSerial)
	return data, err
}

// Get info from this webxdc message.
func (self *Rpc) GetWebxdcInfo(accountId AccountId, msgId MsgId) (*WebxdcMsgInfo, error) {
	var info WebxdcMsgInfo
	err := self.Transport.CallResult(self.Context, &info, "get_webxdc_info", accountId, msgId)
	return &info, err
}

// Get blob encoded as base64 from a webxdc message
//
// path is the path of the file within webxdc archive
func (self *Rpc) GetWebxdcBlob(accountId AccountId, msgId MsgId, path string) (string, error) {
	var data string
	err := self.Transport.CallResult(self.Context, &data, "get_webxdc_blob", accountId, msgId, path)
	return data, err
}

// TODO: get_http_response

// Forward messages to another chat.
//
// All types of messages can be forwarded,
// however, they will be flagged as such (dc_msg_is_forwarded() is set).
//
// Original sender, info-state and webxdc updates are not forwarded on purpose.
func (self *Rpc) ForwardMessages(accountId AccountId, msgIds []MsgId, chatId ChatId) error {
	return self.Transport.Call(self.Context, "forward_messages", accountId, msgIds, chatId)
}

// Resend messages and make information available for newly added chat members.
// Resending sends out the original message, however, recipients and webxdc-status may differ.
// Clients that already have the original message can still ignore the resent message as
// they have tracked the state by dedicated updates.
//
// Some messages cannot be resent, eg. info-messages, drafts, already pending messages or messages that are not sent by SELF.
//
// msgIds all message IDs that should be resend. All messages must belong to the same chat.
func (self *Rpc) ResendMessages(accountId AccountId, msgIds []MsgId) error {
	return self.Transport.Call(self.Context, "resend_messages", accountId, msgIds)
}

func (self *Rpc) SendSticker(accountId AccountId, chatId ChatId, path string) (MsgId, error) {
	var id MsgId
	err := self.Transport.CallResult(self.Context, &id, "send_sticker", accountId, chatId, path)
	return id, err
}

// Send a reaction to message.
//
// Reaction is a string of emojis separated by spaces. Reaction to a
// single message can be sent multiple times. The last reaction
// received overrides all previously received reactions. It is
// possible to remove all reactions by sending an empty string.
func (self *Rpc) SendReaction(accountId AccountId, msgId MsgId, reaction ...string) (MsgId, error) {
	var id MsgId
	err := self.Transport.CallResult(self.Context, &id, "send_reaction", accountId, msgId, reaction)
	return id, err
}

// Returns reactions to the message.
func (self *Rpc) GetMessageReactions(accountId AccountId, msgId MsgId) (option.Option[Reactions], error) {
	var reactions option.Option[Reactions]
	err := self.Transport.CallResult(self.Context, &reactions, "get_message_reactions", accountId, msgId)
	return reactions, err
}

// Send a message and return the resulting Message instance.
func (self *Rpc) SendMsg(accountId AccountId, chatId ChatId, msgData MsgData) (MsgId, error) {
	var id MsgId
	err := self.Transport.CallResult(self.Context, &id, "send_msg", accountId, chatId, msgData)
	return id, err
}

// Checks if messages can be sent to a given chat.
func (self *Rpc) CanSend(accountId AccountId, chatId ChatId) (bool, error) {
	var canSend bool
	err := self.Transport.CallResult(self.Context, &canSend, "can_send", accountId, chatId)
	return canSend, err
}

// ---------------------------------------------
//           functions for the composer
//    the composer is the message input field
// ---------------------------------------------

func (self *Rpc) RemoveDraft(accountId AccountId, chatId ChatId) error {
	return self.Transport.Call(self.Context, "remove_draft", accountId, chatId)
}

// Get draft for a chat, if any.
func (self *Rpc) GetDraft(accountId AccountId, chatId ChatId) (option.Option[MsgSnapshot], error) {
	var msg option.Option[MsgSnapshot]
	err := self.Transport.CallResult(self.Context, &msg, "get_draft", accountId, chatId)
	return msg, err
}

func (self *Rpc) SendVideoChatInvitation(accountId AccountId, chatId ChatId) (MsgId, error) {
	var id MsgId
	err := self.Transport.CallResult(self.Context, &id, "send_videochat_invitation", accountId, chatId)
	return id, err
}

// ---------------------------------------------
//           misc prototyping functions
//       that might get removed later again
// ---------------------------------------------

// TODO: misc_get_sticker_folder()
// TODO: misc_save_sticker()
// TODO: misc_get_stickers()

// Send a text message and return the resulting Message instance.
func (self *Rpc) MiscSendTextMessage(accountId AccountId, chatId ChatId, text string) (MsgId, error) {
	var id MsgId
	err := self.Transport.CallResult(self.Context, &id, "misc_send_text_message", accountId, chatId, text)
	return id, err
}

// TODO: misc_send_msg()

// mimics the old desktop call, will get replaced with something better in the composer rewrite,
// the better version should support:
// - changing viewtype to enable/disable compression
// - keeping same message id as long as attachment does not change for webxdc messages
func (self *Rpc) MiscSetDraft(accountId AccountId, chatId ChatId, text option.Option[string], file option.Option[string], quotedMessageId option.Option[MsgId]) error {
	return self.Transport.Call(self.Context, "misc_set_draft", accountId, chatId, text, file, quotedMessageId)
}
