export interface Reaction {
  id: string
  message_id: string
  user_id: string
  reaction_type: string
  timestamp: string
}

export type ReactionCreateRequest = {
  user_id: string
  reaction_type: string
}

export type ReactionRemoveRequest = {
  user_id: string
  reaction_type: string
}
