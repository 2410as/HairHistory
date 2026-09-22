import { Box, Grid, Text } from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";
import { useTreatments } from "../hooks/useTreatments";
import { TreatmentCard } from "../components/TreatmentCard";
import { SkeletonLoader } from "../components/SkeletonLoader";
import { PageHeader } from "../components/PageHeader";
import { CARD, COLOR, FONT, PRIMARY_BUTTON } from "../design";

interface Props {
  title?: string;
  description?: string;
}

const NewTreatmentButton = () => {
  const navigate = useNavigate();

  return (
    <Box
      as="button"
      {...PRIMARY_BUTTON}
      fontSize="14px"
      px="26px"
      py="13px"
      flexShrink={0}
      onClick={() => navigate("/treatment/new")}
    >
      施術を追加
    </Box>
  );
};

export const TreatmentList = ({
  title = "施術履歴",
  description = "これまでの施術を新しい順に表示しています。",
}: Props) => {
  const { data: treatments, isLoading, error } = useTreatments();

  return (
    <Box>
      <PageHeader title={title} description={description} action={<NewTreatmentButton />} />

      {isLoading && <SkeletonLoader />}

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
            {error.message || "しばらくしてから再度お試しください。"}
          </Text>
        </Box>
      )}

      {!isLoading && !error && treatments.length === 0 && (
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
            まだ記録がありません
          </Text>
          <Text fontFamily={FONT} fontSize="14px" color={COLOR.muted} mt="12px" lineHeight="1.8">
            最初の施術を登録して、髪のストーリーを記録しはじめましょう。
          </Text>
          <Box display="flex" justifyContent="center" mt="28px">
            <NewTreatmentButton />
          </Box>
        </Box>
      )}

      {!isLoading && !error && treatments.length > 0 && (
        <Grid
          templateColumns={{ base: "1fr", md: "repeat(2, 1fr)", lg: "repeat(3, 1fr)" }}
          gap={{ base: "16px", md: "24px" }}
        >
          {treatments.map((treatment) => (
            <TreatmentCard key={treatment.id} treatment={treatment} />
          ))}
        </Grid>
      )}
    </Box>
  );
};
