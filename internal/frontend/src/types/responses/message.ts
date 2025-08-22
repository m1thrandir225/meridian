export type CreateMessageRequest = {
  content_text: string
  is_integration_message?: boolean
  parent_message_id?: string
}
