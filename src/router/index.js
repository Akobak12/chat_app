import { createRouter, createWebHistory } from 'vue-router'
import ChatView from '../views/ChatView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import Homepage from '../views/HomepageView.vue'
import Settings from '../views/SettingsView.vue'
import ServerView from '../views/ServerView.vue'

const routes = [
  {
    path: '/',
    name: 'default',
    redirect: 'login'
  },
  {
    path: '/chat/:id',
    name: 'chat',
    component: ChatView,
    meta: {
      requiresAuth: true
    },
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterView,
  },
  {
    path: '/homepage',
    name: 'homepage',
    component: Homepage,
    meta: {
      requiresAuth: true
    },
  },
  {
    path: '/settings',
    name: 'settings',
    component: Settings,
    meta: {
      requiresAuth: true
    },
  },
  {
    path: '/servers',
    name: 'servers',
    component: ServerView,
    meta: {
      requiresAuth: true
    },
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    console.log(localStorage.getItem("isAuth"));
    if (localStorage.getItem("isAuth") === "false" || localStorage.getItem("isAuth") === null) {
      next({ name: 'login' });
    } else {
      next();
    }
  } else {
    next();
  }
});


export default router
