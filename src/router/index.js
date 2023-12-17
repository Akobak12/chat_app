import { createRouter, createWebHistory } from 'vue-router'
import ChatView from '../views/ChatView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import Homepage from '../views/HomepageView.vue'

const routes = [
  {
    path: '/',
    name: 'default',
    redirect: 'login'
  },
  {
    path: '/chat',
    name: 'chat',
    component: ChatView,
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: {
      disableIfLoggedIn: true
    },
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterView,
    meta: {
      disableIfLoggedIn: true
    },
  },
  {
    path: '/homepage',
    name: 'homepage',
    component: Homepage,
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})



router.beforeEach((to, from, next) => {
  if (localStorage.getItem('isAuth')) {
    if (to.name === 'login' || to.name === 'register') {
      next({ name: 'homepage' });
    } else {
      next();
    }
  } else {
    if (!to.meta.disableIfLoggedIn) {
      next({ name: 'login' });
    } else {
      next();
    }
  }
});

export default router
