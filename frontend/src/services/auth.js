import axios from 'axios';

export async function getUser(opts = {}) {
  return await axios.get("/api/user", opts)
}

export async function logOut(opts = {}) {
  return await axios.get("/api/logout", opts)
}

export async function logIn(username,password,opts = {}) {
    const formData = new FormData()
    formData.append('username', username)
    formData.append('password', password)
    return await axios.post('/api/login',formData)
}