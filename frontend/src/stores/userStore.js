import { defineStore } from 'pinia'
import { ref, reactive,computed } from 'vue'
import { useRouter } from 'vue-router'
import { getUser, logOut, logIn } from '@/services/auth'

export const useUserStore = defineStore('data', () => {
  const router = useRouter()

  const loggedIn = ref(false)
  const user = reactive({})
    
  const isLoggedIn = computed(() => loggedIn.value)

  function setLoggedIn(value) {
    loggedIn.value = value
  }

  async function fetchUser() {
    try {
      const response = await getUser({timeout: 2000})
      Object.assign(user, response.data)
      loggedIn.value = true;
    }
    catch(err) {
      const error = err.response?.data?.message || err.message;
      console.error(`API error: ${error}`);
    }
  }
  
  async function logoutUser() {
    try {
      await logOut()
      loggedIn.value = false;
      router.push('/login')
    } catch (err) {
      const error = err.response?.data?.message || err.message;
      console.error(`API error: ${error}`);
      alert('Error logging out')
    }
  }

async function loginUser(username,password) {
  try {
    await logIn(username,password)
    loggedIn.value = true;
    router.push('/')
   } catch (err) {
    const error = err.response?.data?.message || err.message;
    const status = err.response?.status;

    console.error(`API error: ${error}`);
    if (status === 401) {
      alert("Invalid credentials")
    } else { 
     alert(`Error logging in: ${error}`)
    }
  }
}
  return { loggedIn,user, isLoggedIn, setLoggedIn, fetchUser, logoutUser, loginUser}
});