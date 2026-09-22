export const FONT =
  "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif";

export const COLOR = {
  black: "#000000",
  text: "#000000",
  muted: "#666666",
  faint: "#999999",
  border: "#e5e5e5",
  borderStrong: "#dcdcdc",
  surface: "#f5f5f5",
  white: "#ffffff",
};

export const PAGE_MAX_W = "1160px";

export const PAGE_PX = { base: "24px", md: "40px" };

export const CARD = {
  bg: COLOR.white,
  border: `1px solid ${COLOR.border}`,
  borderRadius: "14px",
};

export const PRIMARY_BUTTON = {
  bg: COLOR.black,
  color: COLOR.white,
  fontFamily: FONT,
  fontWeight: "600",
  borderRadius: "999px",
  border: "none",
  cursor: "pointer",
  transition: "all 200ms ease",
  _hover: { bg: "#333333", transform: "translateY(-1px)", boxShadow: "0 8px 20px rgba(0,0,0,0.18)" },
  _active: { transform: "translateY(0)", boxShadow: "none" },
};

export const SECONDARY_BUTTON = {
  bg: COLOR.white,
  color: COLOR.black,
  fontFamily: FONT,
  fontWeight: "600",
  borderRadius: "999px",
  border: `1px solid ${COLOR.borderStrong}`,
  cursor: "pointer",
  transition: "all 200ms ease",
  _hover: { bg: COLOR.surface, borderColor: "#c4c4c4" },
  _active: { bg: "#ececec" },
};
