import { createRouter, createWebHistory } from "vue-router";
import BrandAnalytics from "../views/brand/Analytics.vue";
import BrandCampaignEditor from "../views/brand/CampaignEditorVant.vue";
import BrandCampaignPageDesigner from "../views/brand/CampaignPageDesigner.vue";
import BrandCampaigns from "../views/brand/Campaigns.vue";
import BrandDashboard from "../views/brand/Dashboard.vue";
import BrandDistributorApproval from "../views/brand/DistributorApproval.vue";
import BrandDistributorLevelRewards from "../views/brand/DistributorLevelRewards.vue";
import BrandDistributors from "../views/brand/Distributors.vue";
import BrandLogin from "../views/brand/Login.vue";
import BrandMaterials from "../views/brand/Materials.vue";
import BrandMemberDetail from "../views/brand/MemberDetail.vue";
import BrandMembers from "../views/brand/Members.vue";
import BrandOrders from "../views/brand/Orders.vue";
import PosterRecords from "../views/brand/PosterRecords.vue";
import BrandPromoters from "../views/brand/Promoters.vue";
import BrandSettings from "../views/brand/Settings.vue";
import VerificationRecords from "../views/brand/VerificationRecords.vue";
import CampaignDetail from "../views/CampaignDetail.vue";
import CampaignForm from "../views/CampaignForm.vue";
import CampaignList from "../views/CampaignList.vue";
import FeedbackCenter from "../views/FeedbackCenter.vue";
import DistributorApply from "../views/distributor/DistributorApply.vue";
import DistributorCenter from "../views/distributor/DistributorCenter.vue";
import DistributorPromotion from "../views/distributor/DistributorPromotion.vue";
import DistributorRewards from "../views/distributor/DistributorRewards.vue";
import DistributorSubordinates from "../views/distributor/DistributorSubordinates.vue";
import DistributorLogin from "../views/distributor/Login.vue";
import MyOrders from "../views/MyOrders.vue";
import OrderVerificationView from "../views/order/OrderVerification.vue";
import PaymentQrcode from "../views/PaymentQrcode.vue";
import PosterGeneratorView from "../views/poster/PosterGenerator.vue";
import Success from "../views/Success.vue";

const routes = [
	{
		path: "/",
		name: "Home",
		component: CampaignList,
	},
	{
		path: "/campaign/:id",
		name: "CampaignDetail",
		component: CampaignDetail,
	},
	{
		path: "/campaign/:id/form",
		name: "CampaignForm",
		component: CampaignForm,
	},
	{
		path: "/orders",
		name: "MyOrders",
		component: MyOrders,
	},
	{
		path: "/feedback",
		name: "FeedbackCenter",
		component: FeedbackCenter,
		meta: { requiresAuth: true },
	},
	{
		path: "/success",
		name: "Success",
		component: Success,
	},
	{
		path: "/verify",
		name: "OrderVerify",
		component: OrderVerificationView,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/campaign/:id/payment",
		name: "PaymentQrcode",
		component: PaymentQrcode,
	},
	{
		path: "/brand/login",
		name: "BrandLogin",
		component: BrandLogin,
	},
	{
		path: "/brand",
		redirect: "/brand/dashboard",
	},
	{
		path: "/brand/dashboard",
		name: "BrandDashboard",
		component: BrandDashboard,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/campaigns",
		name: "BrandCampaigns",
		component: BrandCampaigns,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/campaigns/create",
		name: "BrandCampaignCreate",
		component: BrandCampaignEditor,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/campaigns/edit/:id",
		name: "BrandCampaignEdit",
		component: BrandCampaignEditor,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/campaigns/:id/page-design",
		name: "BrandCampaignPageDesigner",
		component: BrandCampaignPageDesigner,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/materials",
		name: "BrandMaterials",
		component: BrandMaterials,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/orders",
		name: "BrandOrders",
		component: BrandOrders,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/verification-records",
		name: "VerificationRecords",
		component: VerificationRecords,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/poster-records",
		name: "PosterRecords",
		component: PosterRecords,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/promoters",
		name: "BrandPromoters",
		component: BrandPromoters,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/analytics",
		name: "BrandAnalytics",
		component: BrandAnalytics,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/settings",
		name: "BrandSettings",
		component: BrandSettings,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/members",
		name: "BrandMembers",
		component: BrandMembers,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/members/:id",
		name: "BrandMemberDetail",
		component: BrandMemberDetail,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/distributor-approval",
		name: "BrandDistributorApproval",
		component: BrandDistributorApproval,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/distributors",
		name: "BrandDistributors",
		component: BrandDistributors,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/brand/distributor-level-rewards",
		name: "BrandDistributorLevelRewards",
		component: BrandDistributorLevelRewards,
		meta: { requiresAuth: true, role: "brand_admin" },
	},
	{
		path: "/distributor/login",
		name: "DistributorLogin",
		component: DistributorLogin,
	},
	{
		path: "/distributor",
		name: "DistributorCenter",
		component: DistributorCenter,
		meta: { requiresAuth: true, role: "distributor" },
	},
	{
		path: "/distributor/apply",
		name: "DistributorApply",
		component: DistributorApply,
		meta: { requiresAuth: true, role: "distributor" },
	},
	{
		path: "/distributor/promotion",
		name: "DistributorPromotion",
		component: DistributorPromotion,
		meta: { requiresAuth: true, role: "distributor" },
	},
	{
		path: "/distributor/rewards",
		name: "DistributorRewards",
		component: DistributorRewards,
		meta: { requiresAuth: true, role: "distributor" },
	},
	{
		path: "/distributor/subordinates",
		name: "DistributorSubordinates",
		component: DistributorSubordinates,
		meta: { requiresAuth: true, role: "distributor" },
	},
	{
		path: "/poster-generator/:id",
		name: "PosterGenerator",
		component: PosterGeneratorView,
		meta: { requiresAuth: true },
	},
];

const router = createRouter({
	history: createWebHistory(),
	routes,
});

router.beforeEach((to, _from, next) => {
	const token = localStorage.getItem("dmh_token");
	const userRole = localStorage.getItem("dmh_user_role");

	if (to.meta.requiresAuth) {
		if (!token) {
			if (to.path.startsWith("/distributor")) {
				next("/distributor/login");
			} else {
				next("/brand/login");
			}
			return;
		}

		if (to.meta.role && userRole !== to.meta.role) {
			if (to.meta.role === "distributor") {
				next("/distributor/login");
			} else if (to.meta.role === "brand_admin") {
				next("/brand/login");
			} else {
				next("/");
			}
			return;
		}
	}

	next();
});

export default router;
