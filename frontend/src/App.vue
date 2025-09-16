<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { onMounted, watch } from 'vue'
import HeaderComponent from './components/HeaderComponent.vue'
import { useUserStore } from '@/stores/userStore'
import router from '@/router'

const userStore = useUserStore()

watch(() => userStore.isLoggedIn, (val) => {
  if (!val) {
    router.push('/login')
  }
})

onMounted(() => {
  userStore.fetchUser();
})

</script>

<template>
  <div class="d-flex vh-100">
    <!-- Sidebar: logo + navigation only -->
    <aside class="bg-light border-end d-flex flex-column" style="width: 130px;">
      <div class="text-center">
        <img src="@/assets/fcm.svg" class="flogo logo" />
      </div>

      <nav class="flex-grow-1">
        <ul class="nav flex-column">
          <li class="nav-item">
            <RouterLink class="nav-link" to="/">Global</RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink class="nav-link" to="/customers">Customers</RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink class="nav-link" to="/summary">Summary</RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink class="nav-link" to="/about">About</RouterLink>
          </li>
        </ul>
      </nav>

      <div class="mt-auto mx-2 mb-2">
        <button
          v-if="userStore.isLoggedIn"
          @click="userStore.logoutUser"
          class="btn btn-outline-primary btn-sm w-100"
        >
          Logout
        </button>
      </div>
    </aside>

    <!-- Main content: header + router view -->
    <main class="flex-grow-1 d-flex flex-column p-3 overflow-auto">
      <!-- Header -->
      <header class="d-flex justify-content-between align-items-center mb-3">
        <HeaderComponent msg="WWN Manager" />
      </header>

      <!-- Main content -->
      <div class="flex-grow-1">
        <RouterView />
      </div>
    </main>
  </div>
</template>

<style scoped>
.flogo {
  /* max-width: 170px; */
  width: 100%;
  height: auto;
  /* padding: 0rem; */
}
aside .nav-link {
  font-size: 1.1rem; /* adjust as needed, e.g., 1.2rem */
}

a {
  text-decoration: none;
  color: hsla(213, 40%, 58%, 1);
  transition: 0.4s;
  font-weight: 500;
  font-variant: all-small-caps;
}
</style>
