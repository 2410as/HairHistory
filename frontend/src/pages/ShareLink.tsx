import { Box, Text } from "@chakra-ui/react";
import { useState } from "react";
import { QRCodeCanvas } from "qrcode.react";
import { ActionButton } from "../components/ActionButton";
import { PageHeader } from "../components/PageHeader";
import { CARD, COLOR, FONT, PRIMARY_BUTTON, SECONDARY_BUTTON } from "../design";
import { buildShareUrl, useShares } from "../hooks/useShares";
import { toErrorMessage } from "../lib/apiClient";
import { formatExpiresAt } from "../lib/format";

export const ShareLink = () => {
  const { shares, isLoading, error, isCreating, actionError, create, revoke } = useShares();
  const [activeToken, setActiveToken] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const selectedToken = activeToken ?? shares[0]?.token ?? null;
  const shareUrl = selectedToken ? buildShareUrl(selectedToken) : null;

  const handleCreate = async () => {
    const created = await create();
    if (created) setActiveToken(created.token);
  };

  const handleCopyLink = async () => {
    if (!shareUrl) return;
    try {
      await navigator.clipboard.writeText(shareUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      setCopied(false);
    }
  };

  const handleDownloadQR = () => {
    const qrCanvas = document.querySelector("canvas");
    if (!qrCanvas) return;
    const link = document.createElement("a");
    link.href = qrCanvas.toDataURL("image/png");
    link.download = "hairhistory-qr.png";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const handleShare = async () => {
    if (!shareUrl) return;
    if (navigator.share) {
      try {
        await navigator.share({
          title: "HairHistory",
          text: "私の施術履歴です",
          url: shareUrl,
        });
        return;
      } catch {
        return;
      }
    }
    void handleCopyLink();
  };

  const handleRevoke = async (token: string) => {
    const revoked = await revoke(token);
    if (revoked && activeToken === token) setActiveToken(null);
  };

  return (
    <Box>
      <PageHeader
        title="共有"
        description="リンクひとつで、これまでの施術履歴をスタイリストに渡せます。"
        action={
          <ActionButton
            fontSize="14px"
            px="26px"
            py="13px"
            flexShrink={0}
            disabled={isCreating}
            onClick={handleCreate}
          >
            {isCreating ? "発行中…" : "共有リンクを発行"}
          </ActionButton>
        }
      />

      <Box maxW="560px">
        {isLoading && (
          <Box
            bg={CARD.bg}
            border={CARD.border}
            borderRadius={CARD.borderRadius}
            px={{ base: "24px", md: "40px" }}
            py={{ base: "48px", md: "64px" }}
            textAlign="center"
          >
            <Text fontFamily={FONT} fontSize="14px" color={COLOR.muted}>
              読み込み中…
            </Text>
          </Box>
        )}

        {!isLoading && error && (
          <Box
            bg={CARD.bg}
            border={CARD.border}
            borderRadius={CARD.borderRadius}
            px={{ base: "24px", md: "40px" }}
            py={{ base: "48px", md: "64px" }}
            textAlign="center"
          >
            <Text fontFamily={FONT} fontSize="16px" fontWeight="600" color={COLOR.black}>
              読み込みに失敗しました
            </Text>
            <Text fontFamily={FONT} fontSize="14px" color={COLOR.muted} mt="10px" lineHeight="1.8">
              {toErrorMessage(error)}
            </Text>
          </Box>
        )}

        {!isLoading && !error && !shareUrl && (
          <Box
            bg={CARD.bg}
            border={CARD.border}
            borderRadius={CARD.borderRadius}
            px={{ base: "24px", md: "40px" }}
            py={{ base: "56px", md: "80px" }}
            textAlign="center"
          >
            <Text
              fontFamily={FONT}
              fontSize={{ base: "18px", md: "20px" }}
              fontWeight="700"
              color={COLOR.black}
            >
              有効な共有リンクがありません
            </Text>
            <Text fontFamily={FONT} fontSize="14px" color={COLOR.muted} mt="12px" lineHeight="1.8">
              リンクを発行すると、7 日間だけ施術履歴を見せられる URL と QR コードが手に入ります。
            </Text>
          </Box>
        )}

        {!isLoading && !error && shareUrl && (
          <>
            <Box
              bg={CARD.bg}
              border={CARD.border}
              borderRadius={CARD.borderRadius}
              p={{ base: "24px", md: "32px" }}
              display="flex"
              flexDirection="column"
              alignItems="center"
            >
              <Box bg={COLOR.white} p="16px" borderRadius="12px" border={`1px solid ${COLOR.border}`}>
                <QRCodeCanvas
                  value={shareUrl}
                  size={180}
                  level="H"
                  bgColor="#ffffff"
                  fgColor="#000000"
                />
              </Box>
              <Text fontFamily={FONT} fontSize="14px" color={COLOR.muted} mt="20px" textAlign="center">
                QR コードを読み取ると施術履歴を開けます
              </Text>
            </Box>

            <Box mt="24px">
              <Text fontFamily={FONT} fontSize="14px" fontWeight="600" color={COLOR.black} mb="10px">
                共有リンク
              </Text>
              <Box
                bg={COLOR.surface}
                border={`1px solid ${COLOR.border}`}
                borderRadius="10px"
                px="16px"
                py="14px"
                userSelect="all"
                wordBreak="break-all"
              >
                <Text
                  fontFamily={FONT}
                  fontSize="14px"
                  color={COLOR.muted}
                  lineHeight="1.6"
                  data-testid="share-url"
                >
                  {shareUrl}
                </Text>
              </Box>
            </Box>

            <Box
              display="flex"
              flexDirection={{ base: "column", sm: "row" }}
              flexWrap="wrap"
              gap="12px"
              mt="24px"
            >
              <Box as="button" {...PRIMARY_BUTTON} fontSize="15px" px="32px" py="14px" onClick={handleShare}>
                共有する
              </Box>
              <Box as="button" {...SECONDARY_BUTTON} fontSize="15px" px="28px" py="14px" onClick={handleCopyLink}>
                {copied ? "コピーしました" : "リンクをコピー"}
              </Box>
              <Box as="button" {...SECONDARY_BUTTON} fontSize="15px" px="28px" py="14px" onClick={handleDownloadQR}>
                QR を保存
              </Box>
            </Box>
          </>
        )}

        {actionError && (
          <Text fontFamily={FONT} fontSize="13px" color={COLOR.muted} mt="16px" lineHeight="1.7">
            {toErrorMessage(actionError)}
          </Text>
        )}

        {!isLoading && !error && shares.length > 0 && (
          <Box mt="40px">
            <Text fontFamily={FONT} fontSize="14px" fontWeight="600" color={COLOR.black} mb="12px">
              有効なリンク（{shares.length}）
            </Text>
            <Box display="flex" flexDirection="column" gap="12px">
              {shares.map((share) => (
                <Box
                  key={share.id}
                  bg={CARD.bg}
                  border={CARD.border}
                  borderRadius={CARD.borderRadius}
                  p="16px"
                  display="flex"
                  flexDirection={{ base: "column", sm: "row" }}
                  gap="12px"
                  justifyContent="space-between"
                  alignItems={{ base: "flex-start", sm: "center" }}
                >
                  <Box minWidth="0" flex="1">
                    <Text
                      fontFamily={FONT}
                      fontSize="13px"
                      color={share.token === selectedToken ? COLOR.black : COLOR.muted}
                      fontWeight={share.token === selectedToken ? "600" : "500"}
                      wordBreak="break-all"
                      lineHeight="1.6"
                    >
                      {buildShareUrl(share.token)}
                    </Text>
                    <Text fontFamily={FONT} fontSize="12px" color={COLOR.faint} mt="4px">
                      有効期限 {formatExpiresAt(share.expiresAt)}
                    </Text>
                  </Box>
                  <Box display="flex" gap="8px" flexShrink={0}>
                    <Box
                      as="button"
                      {...SECONDARY_BUTTON}
                      fontSize="13px"
                      px="18px"
                      py="9px"
                      onClick={() => setActiveToken(share.token)}
                    >
                      表示
                    </Box>
                    <Box
                      as="button"
                      {...SECONDARY_BUTTON}
                      fontSize="13px"
                      px="18px"
                      py="9px"
                      onClick={() => handleRevoke(share.token)}
                    >
                      無効化
                    </Box>
                  </Box>
                </Box>
              ))}
            </Box>
          </Box>
        )}
      </Box>
    </Box>
  );
};
