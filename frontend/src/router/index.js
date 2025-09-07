import { createRouter, createWebHistory } from 'vue-router'
import FCManager  from '../views/FCManager.vue'
import GlobalFCManager  from '../views/GlobalFCManager.vue'
import About  from '../views/About.vue'
import Summary  from '../views/Summary.vue'
import Login from '../views/Login.vue'
import { useUserStore } from '@/stores/userStore'
import { useApiStore } from '@/stores/apiStore'


const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      alias: ['/index.html'],
      name: 'global',
      component: GlobalFCManager,
      meta: { requiresData: true },
    },
    {
      path: '/customers',
      alias: ['/customers.html'],
      name: 'customers',
      component: FCManager,
      meta: { requiresData: true },
    },
    {
      path: '/summary',
      alias: ['/summary.html'],
      name: 'summary',
      component: Summary,
    },
    {
      path: '/about',
      alias: ['/about.html'],
      name: 'about',
      component: About 
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
        const userStore = useUserStore()
        if (userStore.isLoggedIn) {
          next('/') // redirect if already logged in
        } else {
          next()
        }
      }
    }
  ]
})

router.beforeEach(async (to, from, next) => {
  if (to.meta.requiresData) {
    useApiStore().init();
  }
  next()
})

export default router
