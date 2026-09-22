import { Box, Heading, Text } from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";
import { FONT } from "../design";

const TAGS = ["施術記録", "かんたん評価", "共有"];

const FEATURES = [
  {
    title: "シンプルで洗練されたデザイン",
    body: "モノトーンで統一した画面構成。余計な装飾を削ぎ落とし、必要な情報だけが自然に目に入ります。はじめて開いた人でも迷わず使えます。",
  },
  {
    title: "施術履歴の管理がスムーズに",
    body: "カット・カラー・トリートメントの記録を日付順に一覧表示。施術の傾向はグラフで振り返れるので、次に何をすべきかがひと目でわかります。",
  },
  {
    title: "共有がかんたんに",
    body: "リンクひとつで、これまでの施術履歴をスタイリストに渡せます。口頭で説明しづらい過去の薬剤や仕上がりも、正確に伝わります。",
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

const ScreenHistory = () => {
  const rows = [
    { date: "09 / 15", label: "カット" },
    { date: "09 / 01", label: "カラー" },
    { date: "08 / 20", label: "トリートメント" },
  ];
  const bars = [38, 56, 30, 70, 48, 62];

  return (
    <>
      <ScreenTitle>施術履歴</ScreenTitle>
      <Box display="flex" flexDirection="column" gap="6px">
        {rows.map((row) => (
          <Box
            key={row.date}
            bg="#f5f5f5"
            borderRadius="8px"
            px="10px"
            py="8px"
            display="flex"
            justifyContent="space-between"
            alignItems="center"
          >
            <Text fontFamily={FONT} fontSize="9px" fontWeight="600" color="#000000">
              {row.label}
            </Text>
            <Text fontFamily={FONT} fontSize="9px" color="#666666">
              {row.date}
            </Text>
          </Box>
        ))}
      </Box>
      <Box
        borderTop="1px solid #e5e5e5"
        pt="10px"
        mt="2px"
        flex="1"
        display="flex"
        flexDirection="column"
      >
        <Text fontFamily={FONT} fontSize="9px" color="#666666" mb="8px">
          来店ペース
        </Text>
        <Box display="flex" alignItems="flex-end" gap="6px" flex="1" pb="4px">
          {bars.map((height, index) => (
            <Box
              key={index}
              flex="1"
              height={`${height}%`}
              borderRadius="3px"
              bg={index === 3 ? "#000000" : "#d6d6d6"}
            />
          ))}
        </Box>
      </Box>
    </>
  );
};

const StarRow = ({ label, score }: { label: string; score: number }) => (
  <Box display="flex" justifyContent="space-between" alignItems="center">
    <Text fontFamily={FONT} fontSize="9px" color="#666666">
      {label}
    </Text>
    <Box display="flex" gap="2px">
      {[1, 2, 3, 4, 5].map((star) => (
        <Text key={star} fontSize="12px" lineHeight="1" color={star <= score ? "#000000" : "#dcdcdc"}>
          ★
        </Text>
      ))}
    </Box>
  </Box>
);

const ScreenRating = () => (
  <>
    <ScreenTitle>評価を入力</ScreenTitle>
    <Box bg="#f5f5f5" borderRadius="10px" px="12px" py="12px" display="flex" flexDirection="column" gap="10px">
      <StarRow label="仕上がり" score={5} />
      <StarRow label="色もち" score={4} />
      <StarRow label="扱いやすさ" score={3} />
    </Box>
    <Box border="1px solid #e5e5e5" borderRadius="10px" px="10px" py="10px" flex="1">
      <Text fontFamily={FONT} fontSize="9px" color="#999999" lineHeight="1.6">
        メモを残す
      </Text>
      <Box mt="8px" display="flex" flexDirection="column" gap="5px">
        <Box height="5px" bg="#ececec" borderRadius="3px" />
        <Box height="5px" bg="#ececec" borderRadius="3px" />
        <Box height="5px" bg="#ececec" borderRadius="3px" width="60%" />
      </Box>
    </Box>
    <Box bg="#000000" borderRadius="8px" py="9px" textAlign="center">
      <Text fontFamily={FONT} fontSize="10px" fontWeight="600" color="#ffffff">
        保存する
      </Text>
    </Box>
  </>
);

const DetailRow = ({ label, value }: { label: string; value: string }) => (
  <Box display="flex" justifyContent="space-between" alignItems="center" py="7px" borderBottom="1px solid #efefef">
    <Text fontFamily={FONT} fontSize="9px" color="#666666">
      {label}
    </Text>
    <Text fontFamily={FONT} fontSize="9px" fontWeight="600" color="#000000">
      {value}
    </Text>
  </Box>
);

const ScreenDetail = () => (
  <>
    <ScreenTitle>施術詳細</ScreenTitle>
    <Box bg="#f5f5f5" borderRadius="10px" px="12px" py="12px">
      <Text fontFamily={FONT} fontSize="13px" fontWeight="700" color="#000000">
        カラー
      </Text>
      <Text fontFamily={FONT} fontSize="9px" color="#666666" mt="2px">
        2026 / 09 / 01
      </Text>
    </Box>
    <Box>
      <DetailRow label="サロン" value="HAIR STUDIO" />
      <DetailRow label="担当" value="Yuki" />
      <DetailRow label="薬剤" value="8 トーン" />
      <DetailRow label="料金" value="¥ 9,800" />
    </Box>
    <Box border="1px solid #e5e5e5" borderRadius="10px" px="10px" py="10px" flex="1">
      <Text fontFamily={FONT} fontSize="9px" color="#999999">
        メモ
      </Text>
      <Text fontFamily={FONT} fontSize="9px" color="#666666" lineHeight="1.7" mt="5px">
        明るめの仕上がり。 次回は 少し暗めに 調整する。
      </Text>
    </Box>
  </>
);

const SCREENS = [<ScreenHistory key="history" />, <ScreenRating key="rating" />, <ScreenDetail key="detail" />];

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
          あなたの髪のストーリーを、記録して、評価して、共有する。
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
