import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../views/home/HomePage.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/home' },
    { path: '/home', component: HomePage }
  ]
})
