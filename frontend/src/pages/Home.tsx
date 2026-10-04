import { Box, Heading, Text } from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";
import { FONT } from "../design";

const TAGS = ["施術履歴", "施術の登録", "リンク共有"];

const FEATURES = [
  {
    title: "シンプルで洗練されたデザイン",
    body: "モノトーンで統一した画面構成。余計な装飾を削ぎ落とし、必要な情報だけが自然に目に入ります。はじめて開いた人でも迷わず使えます。",
  },
  {
    title: "施術履歴をまとめて見返せる",
    body: "日付・施術項目・美容院名・メモ・料金を記録しておくと、新しい順にカードで一覧表示されます。カット・カラー・縮毛矯正のように複数の施術をまとめて1件に登録でき、あとから編集・削除もできます。",
  },
  {
    title: "リンクひとつで共有できる",
    body: "有効期限つきの共有リンクを発行して、これまでの施術履歴をスタイリストに渡せます。リンクのコピーと QR コードの保存に対応していて、不要になったリンクはいつでも無効化できます。",
  },
];

const PhoneFrame = ({ children }: { children: React.ReactNode }) => (
  <Box
    width="100%"
    maxW="248px"
    margin="0 auto"
    border="12px solid #000000"
    borderRadius="44px"
    bg="#000000"
    boxShadow="0 18px 40px rgba(0,0,0,0.18)"
    overflow="hidden"
  >
    <Box
      position="relative"
      width="100%"
      height={{ base: "440px", md: "470px" }}
      bg="#ffffff"
      borderRadius="32px"
      overflow="hidden"
    >
      <Box
        position="absolute"
        top="0"
        left="50%"
        transform="translateX(-50%)"
        width="96px"
        height="22px"
        bg="#000000"
        borderRadius="0 0 14px 14px"
        zIndex="2"
      />
      <Box height="100%" pt="34px" px="14px" pb="14px" display="flex" flexDirection="column" gap="10px">
        {children}
      </Box>
    </Box>
  </Box>
);

const ScreenTitle = ({ children }: { children: React.ReactNode }) => (
  <Text fontFamily={FONT} fontSize="12px" fontWeight="700" color="#000000" letterSpacing="0.02em">
    {children}
  </Text>
);

const ServiceChip = ({ children }: { children: React.ReactNode }) => (
  <Box bg="#ffffff" border="1px solid #e5e5e5" borderRadius="999px" px="7px" py="3px">
    <Text fontFamily={FONT} fontSize="8px" fontWeight="600" color="#000000" lineHeight="1.2">
      {children}
    </Text>
  </Box>
);

const HistoryCard = ({
  services,
  salon,
  date,
  cost,
}: {
  services: string[];
  salon: string;
  date: string;
  cost: string;
}) => (
  <Box bg="#f5f5f5" borderRadius="10px" px="10px" py="9px">
    <Box display="flex" flexWrap="wrap" gap="4px">
      {services.map((service) => (
        <ServiceChip key={service}>{service}</ServiceChip>
      ))}
    </Box>
    <Text fontFamily={FONT} fontSize="10px" fontWeight="700" color="#000000" mt="7px">
      {salon}
    </Text>
    <Box display="flex" justifyContent="space-between" alignItems="center" mt="4px">
      <Text fontFamily={FONT} fontSize="8px" color="#999999">
        {date}
      </Text>
      <Text fontFamily={FONT} fontSize="9px" fontWeight="600" color="#000000">
        {cost}
      </Text>
    </Box>
  </Box>
);

const ScreenHistory = () => (
  <>
    <ScreenTitle>施術履歴</ScreenTitle>
    <Box display="flex" flexDirection="column" gap="7px">
      <HistoryCard
        services={["カット", "カラー"]}
        salon="HAIR STUDIO"
        date="2026 / 09 / 15"
        cost="¥12,000"
      />
      <HistoryCard
        services={["縮毛矯正"]}
        salon="Hair Salon ABC"
        date="2026 / 09 / 01"
        cost="¥18,000"
      />
      <HistoryCard
        services={["トリートメント"]}
        salon="HAIR STUDIO"
        date="2026 / 08 / 20"
        cost="¥6,500"
      />
      <HistoryCard
        services={["カット", "ヘッドスパ"]}
        salon="Hair Salon ABC"
        date="2026 / 07 / 28"
        cost="¥9,800"
      />
    </Box>
  </>
);

const FormField = ({ label, value }: { label: string; value: string }) => (
  <Box>
    <Text fontFamily={FONT} fontSize="8px" fontWeight="600" color="#666666" mb="3px">
      {label}
    </Text>
    <Box border="1px solid #e5e5e5" borderRadius="7px" px="8px" py="6px">
      <Text fontFamily={FONT} fontSize="9px" color="#000000" lineHeight="1.3">
        {value}
      </Text>
    </Box>
  </Box>
);

const ScreenForm = () => (
  <>
    <ScreenTitle>施術を追加</ScreenTitle>
    <Box display="flex" flexDirection="column" gap="7px" flex="1">
      <FormField label="日付" value="2026 / 09 / 15" />
      <Box>
        <Text fontFamily={FONT} fontSize="8px" fontWeight="600" color="#666666" mb="3px">
          施術項目
        </Text>
        <Box display="flex" flexWrap="wrap" gap="4px">
          <Box bg="#000000" borderRadius="999px" px="8px" py="4px">
            <Text fontFamily={FONT} fontSize="8px" fontWeight="600" color="#ffffff" lineHeight="1.2">
              カット
            </Text>
          </Box>
          <Box bg="#000000" borderRadius="999px" px="8px" py="4px">
            <Text fontFamily={FONT} fontSize="8px" fontWeight="600" color="#ffffff" lineHeight="1.2">
              カラー
            </Text>
          </Box>
          <ServiceChip>パーマ</ServiceChip>
          <ServiceChip>ヘッドスパ</ServiceChip>
        </Box>
      </Box>
      <FormField label="美容院名" value="HAIR STUDIO" />
      <FormField label="料金" value="12000" />
      <Box flex="1" display="flex" flexDirection="column">
        <Text fontFamily={FONT} fontSize="8px" fontWeight="600" color="#666666" mb="3px">
          メモ
        </Text>
        <Box border="1px solid #e5e5e5" borderRadius="7px" px="8px" py="6px" flex="1">
          <Text fontFamily={FONT} fontSize="8px" color="#666666" lineHeight="1.6">
            8 トーンのアッシュ。 次回は 少し暗めに 調整する。
          </Text>
        </Box>
      </Box>
    </Box>
    <Box bg="#000000" borderRadius="8px" py="9px" textAlign="center">
      <Text fontFamily={FONT} fontSize="10px" fontWeight="600" color="#ffffff">
        保存する
      </Text>
    </Box>
  </>
);

const QR_PATTERN = [
  "1110111011101",
  "1000101010001",
  "1011100101101",
  "1010010110101",
  "1110101011101",
  "0000110100000",
  "1101011101011",
  "0100100010110",
  "1011101011101",
  "0001010110010",
  "1110100011101",
  "1000110101001",
  "1110101110111",
];

const QRCodeMock = () => (
  <Box
    bg="#ffffff"
    border="1px solid #e5e5e5"
    borderRadius="10px"
    p="10px"
    display="flex"
    justifyContent="center"
  >
    <Box
      display="grid"
      gridTemplateColumns={`repeat(${QR_PATTERN.length}, 1fr)`}
      width="104px"
      aria-hidden="true"
    >
      {QR_PATTERN.flatMap((row, rowIndex) =>
        row.split("").map((cell, cellIndex) => (
          <Box
            key={`${rowIndex}-${cellIndex}`}
            width="100%"
            paddingTop="100%"
            bg={cell === "1" ? "#000000" : "#ffffff"}
          />
        ))
      )}
    </Box>
  </Box>
);

const ScreenShare = () => (
  <>
    <ScreenTitle>共有リンク</ScreenTitle>
    <Box bg="#f5f5f5" borderRadius="8px" px="10px" py="8px">
      <Text fontFamily={FONT} fontSize="8px" color="#999999">
        共有 URL
      </Text>
      <Text fontFamily={FONT} fontSize="8px" color="#000000" mt="3px" lineHeight="1.4">
        hairhistory.app/shares/ ••••••••••••
      </Text>
    </Box>
    <QRCodeMock />
    <Box display="flex" justifyContent="space-between" alignItems="center" flex="1">
      <Text fontFamily={FONT} fontSize="8px" color="#666666">
        有効期限
      </Text>
      <Text fontFamily={FONT} fontSize="8px" fontWeight="600" color="#000000">
        2026 / 09 / 22 まで
      </Text>
    </Box>
    <Box display="flex" gap="6px">
      <Box bg="#000000" borderRadius="8px" py="9px" textAlign="center" flex="1">
        <Text fontFamily={FONT} fontSize="9px" fontWeight="600" color="#ffffff">
          リンクをコピー
        </Text>
      </Box>
      <Box
        bg="#ffffff"
        border="1px solid #dcdcdc"
        borderRadius="8px"
        py="9px"
        textAlign="center"
        flex="1"
      >
        <Text fontFamily={FONT} fontSize="9px" fontWeight="600" color="#000000">
          QR を保存
        </Text>
      </Box>
    </Box>
  </>
);

const SCREENS = [<ScreenHistory key="history" />, <ScreenForm key="form" />, <ScreenShare key="share" />];

export const Home = () => {
  const navigate = useNavigate();

  return (
    <Box width="100%" bg="#ffffff" minHeight="100vh">
      <Box maxW="1160px" margin="0 auto" px={{ base: "24px", md: "40px" }} py={{ base: "64px", md: "110px" }}>
        <Heading
          as="h1"
          textAlign="center"
          fontFamily={FONT}
          fontSize={{ base: "34px", md: "52px" }}
          fontWeight="700"
          lineHeight="1.25"
          letterSpacing="-0.01em"
          color="#000000"
        >
          HairHistory でできること
        </Heading>
        <Text
          textAlign="center"
          fontFamily={FONT}
          fontSize={{ base: "14px", md: "16px" }}
          color="#666666"
          lineHeight="1.8"
          mt={{ base: "16px", md: "20px" }}
        >
          あなたの髪のストーリーを、記録して、共有する。
        </Text>

        <Box
          display="grid"
          gridTemplateColumns={{ base: "1fr", md: "repeat(3, 1fr)" }}
          gap={{ base: "48px", md: "40px" }}
          mt={{ base: "56px", md: "80px" }}
        >
          {TAGS.map((tag, index) => (
            <Box key={tag} display="flex" flexDirection="column" alignItems="center">
              <Box
                bg="#ffffff"
                border="1px solid #dcdcdc"
                borderRadius="999px"
                px="24px"
                py="9px"
                mb={{ base: "24px", md: "32px" }}
              >
                <Text fontFamily={FONT} fontSize={{ base: "13px", md: "14px" }} fontWeight="600" color="#000000">
                  {tag}
                </Text>
              </Box>
              <PhoneFrame>{SCREENS[index]}</PhoneFrame>
            </Box>
          ))}
        </Box>

        <Box
          maxW="780px"
          margin="0 auto"
          mt={{ base: "72px", md: "110px" }}
          display="flex"
          flexDirection="column"
          gap={{ base: "36px", md: "44px" }}
        >
          {FEATURES.map((feature) => (
            <Box key={feature.title} borderTop="1px solid #e5e5e5" pt={{ base: "24px", md: "30px" }}>
              <Heading
                as="h2"
                fontFamily={FONT}
                fontSize={{ base: "19px", md: "23px" }}
                fontWeight="700"
                lineHeight="1.4"
                color="#000000"
                mb={{ base: "10px", md: "14px" }}
              >
                {feature.title}
              </Heading>
              <Text fontFamily={FONT} fontSize={{ base: "14px", md: "16px" }} color="#666666" lineHeight="1.9">
                {feature.body}
              </Text>
            </Box>
          ))}
        </Box>

        <Box display="flex" justifyContent="center" mt={{ base: "64px", md: "96px" }}>
          <Box
            as="button"
            onClick={() => navigate("/dashboard")}
            bg="#000000"
            color="#ffffff"
            fontFamily={FONT}
            fontSize={{ base: "15px", md: "16px" }}
            fontWeight="600"
            borderRadius="999px"
            border="none"
            px={{ base: "44px", md: "56px" }}
            py={{ base: "18px", md: "21px" }}
            cursor="pointer"
            transition="all 200ms ease"
            _hover={{ bg: "#333333", transform: "translateY(-2px)", boxShadow: "0 10px 24px rgba(0,0,0,0.2)" }}
          >
            今すぐ始める
          </Box>
        </Box>
      </Box>
    </Box>
  );
};
