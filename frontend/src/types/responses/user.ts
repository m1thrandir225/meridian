import type { User } from '../models/user'

export type UpdateUserRequest = Partial<Omit<User, 'id'>>

export type UpdateUserResponse = User

export type UpdateUserPasswordRequest = {
  new_password: string
}
