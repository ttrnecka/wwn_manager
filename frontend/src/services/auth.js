import axios from 'axios';
import fcService from "@/services/fcService";

export async function getUser(opts = {}) {
  return await fcService.api.get("/user", opts)
}

export async function logOut(opts = {}) {
  return await fcService.api.get("/logout", opts)
}

export async function logIn(username,password) {
    const formData = new FormData()
    formData.append('username', username)
    formData.append('password', password)
    return await fcService.api.post('/login',formData)
}