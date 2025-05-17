export interface UserLoginType {
  username: string
  password: string
}

export interface UserType {
  username: string
  id: number
  avatar: string
  password?: string
  role?: string
  roleId?: string
}
