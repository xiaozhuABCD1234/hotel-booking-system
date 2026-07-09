import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "Login",
      component: () => import("@/views/auth/LoginView.vue"),
      meta: { guest: true },
    },
    {
      path: "/register",
      name: "Register",
      component: () => import("@/views/auth/RegisterView.vue"),
      meta: { guest: true },
    },
    // ── User Routes ────────────────────────────
    {
      path: "/",
      component: () => import("@/layouts/DefaultLayout.vue"),
      children: [
        {
          path: "",
          name: "Home",
          component: () => import("@/views/user/HomeView.vue"),
        },
        {
          path: "hotel/:id",
          name: "HotelDetail",
          component: () => import("@/views/user/HotelDetailView.vue"),
        },
        {
          path: "booking/:roomId",
          name: "Booking",
          component: () => import("@/views/user/BookingView.vue"),
          meta: { requiresAuth: true },
        },
        {
          path: "orders",
          name: "MyOrders",
          component: () => import("@/views/user/MyOrdersView.vue"),
          meta: { requiresAuth: true },
        },
        {
          path: "profile",
          name: "Profile",
          component: () => import("@/views/user/ProfileView.vue"),
          meta: { requiresAuth: true },
        },
      ],
    },
    // ── Admin Routes ───────────────────────────
    {
      path: "/admin",
      component: () => import("@/layouts/AdminLayout.vue"),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: "",
          name: "Dashboard",
          component: () => import("@/views/admin/DashboardView.vue"),
          meta: { requiresAuth: true, requiresAdmin: true },
        },
        {
          path: "hotels",
          name: "HotelManage",
          component: () => import("@/views/admin/HotelManageView.vue"),
          meta: { requiresAuth: true, requiresAdmin: true },
        },
        {
          path: "rooms",
          name: "RoomManage",
          component: () => import("@/views/admin/RoomManageView.vue"),
          meta: { requiresAuth: true, requiresAdmin: true },
        },
        {
          path: "orders",
          name: "OrderManage",
          component: () => import("@/views/admin/OrderManageView.vue"),
          meta: { requiresAuth: true, requiresAdmin: true },
        },
        {
          path: "users",
          name: "UserManage",
          component: () => import("@/views/admin/UserManageView.vue"),
          meta: { requiresAuth: true, requiresAdmin: true },
        },
        {
          path: "reports",
          name: "Reports",
          component: () => import("@/views/admin/ReportView.vue"),
          meta: { requiresAuth: true, requiresAdmin: true },
        },
      ],
    },
    // ── 404 ────────────────────────────────────
    {
      path: "/:pathMatch(.*)*",
      name: "NotFound",
      component: () => import("@/views/NotFoundView.vue"),
    },
  ],
});

router.beforeEach((to, _from) => {
  const auth = useAuthStore();

  console.log(
    "[Router]",
    "to:",
    to.path,
    "meta:",
    to.meta,
    "isAdmin:",
    auth.isAdmin,
    "user:",
    auth.user,
  );

  // Guest-only routes (login/register): redirect to home if already logged in
  if (to.meta.guest && auth.isLoggedIn) {
    return { name: "Home", replace: true };
  }

  // Auth-required routes
  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return { name: "Login", query: { redirect: to.fullPath }, replace: true };
  }

  // Admin-only routes
  if (to.meta.requiresAdmin && !auth.isAdmin) {
    return { name: "Home", replace: true };
  }
});

export default router;
