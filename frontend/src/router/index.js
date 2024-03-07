import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import UserRegister from '@/views/UserRegister.vue'
import UserLogin from '@/views/UserLogin.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/register',
      name: 'register',
      component: UserRegister
    },
    {
      path: '/login',
      name: 'login',
      component: UserLogin
    }
  ]
})

export default router
