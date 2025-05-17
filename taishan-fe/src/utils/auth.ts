import { useCookies } from 'vue3-cookies'

const TokenKey = 'token'
const Domain = '.cdfsunrise.com'

const { cookies } = useCookies()

export function getToken() {
  return 'token'
  // return cookies.get(TokenKey)
}

export function setToken(token) {
  return cookies.set(TokenKey, token)
}

export function removeToken() {
  return cookies.remove(TokenKey, Domain)
}
