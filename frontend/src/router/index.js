import { createRouter, createWebHistory } from 'vue-router'
import FCManager  from '../views/FCManager.vue'
import GlobalFCManager  from '../views/GlobalFCManager.vue'
import Login from '../views/Login.vue'
import { useDataStore } from '@/stores/dataStore'


const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      alias: ['/index.html'],
      name: 'global',
      component: GlobalFCManager 
    },
    {
      path: '/customers',
      alias: ['/customers.html'],
      name: 'customers',
      component: FCManager 
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    },
    {
    path: '/login',
      component: Login,
      beforeEnter: (to, from, next) => {
        const dataStore = useDataStore()
        if (dataStore.isLoggedIn) {
          next('/') // redirect if already logged in
        } else {
          next()
        }
      }
    }
  ]
})

export default router
