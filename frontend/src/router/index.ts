import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBucketStore } from '@/stores/bucket'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('@/views/SignUpView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/buckets',
      name: 'buckets',
      component: () => import('@/views/BucketSelectView.vue'),
      meta: { requiresAuth: true, requiresBucket: false },
    },
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/files',
      name: 'files',
      component: () => import('@/views/FilesView.vue'),
      meta: { requiresAuth: true, requiresBucket: true },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const bucketStore = useBucketStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login' })
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    // Если пользователь авторизован, перенаправляем на выбор бакета
    if (!bucketStore.selectedBucketId) {
      next({ name: 'buckets' })
    } else {
      next({ name: 'home' })
    }
  } else if (to.meta.requiresBucket && !bucketStore.selectedBucketId) {
    // Если требуется бакет, но он не выбран, перенаправляем на выбор
    next({ name: 'buckets' })
  } else {
    next()
  }
})

export default router
