export const designTokens = {
  color: {
    white: "#FFFFFF",
    textPrimary: "#1A2B3F",
    textSecondary: "#566B84",
    border: "#D8E4EF",
    userBrand: "#2D85A2",
    adminBrand: "#2B7D9A",
    userBgStart: "#F8FBFF",
    userBgEnd: "#EEF4FB",
    adminBgStart: "#F9FBFD",
    adminBgEnd: "#EDF3F8"
  },
  radius: {
    md: "12rpx",
    lg: "24rpx"
  },
  spacing: {
    pageX: "32rpx",
    card: "36rpx",
    gap: "16rpx"
  },
  fontSize: {
    title: "42rpx",
    body: "28rpx",
    meta: "22rpx"
  }
} as const;

export const userThemeTokens = {
  brand: designTokens.color.userBrand,
  pageBackground: `linear-gradient(180deg, ${designTokens.color.userBgStart} 0%, ${designTokens.color.userBgEnd} 100%)`
} as const;

export const adminThemeTokens = {
  brand: designTokens.color.adminBrand,
  pageBackground: `linear-gradient(180deg, ${designTokens.color.adminBgStart} 0%, ${designTokens.color.adminBgEnd} 100%)`
} as const;
