import { Box, Grid, Text } from '@chakra-ui/react';
import { useParams } from 'react-router-dom';
import { PageHeader } from '../components/PageHeader';
import { SkeletonLoader } from '../components/SkeletonLoader';
import { ServiceTags } from '../components/ServiceTags';
import { CARD, COLOR, FONT, PAGE_MAX_W, PAGE_PX } from '../design';
import { usePublicShare } from '../hooks/useShares';
import { ApiError } from '../lib/apiClient';
import { formatTreatedOn } from '../lib/format';
import type { PublicTreatment } from '../types';

const PublicTreatmentCard = ({ treatment }: { treatment: PublicTreatment }) => (
  <Box
    bg={CARD.bg}
    border={CARD.border}
    borderRadius={CARD.borderRadius}
    p={{ base: '20px', md: '24px' }}
    display="flex"
    flexDirection="column"
    minHeight="200px"
  >
    <ServiceTags services={treatment.services} />

    <Text
      fontFamily={FONT}
      fontSize={{ base: '17px', md: '18px' }}
      fontWeight="700"
      color={COLOR.black}
      lineHeight="1.4"
      mt="16px"
    >
      {treatment.salonName || 'サロン名なし'}
    </Text>

    <Text fontFamily={FONT} fontSize="13px" color={COLOR.faint} mt="6px">
      {formatTreatedOn(treatment.treatedOn)}
    </Text>

    {treatment.memo && (
      <Text fontFamily={FONT} fontSize="14px" color={COLOR.muted} lineHeight="1.8" mt="14px">
        {treatment.memo}
      </Text>
    )}

    <Box mt="auto" pt="20px">
      <Text fontFamily={FONT} fontSize="15px" fontWeight="600" color={COLOR.black}>
        {typeof treatment.cost === 'number' ? `¥${treatment.cost.toLocaleString()}` : ''}
      </Text>
    </Box>
  </Box>
);

const EmptyState = ({ title, body }: { title: string; body: string }) => (
  <Box
    bg={CARD.bg}
    border={CARD.border}
    borderRadius={CARD.borderRadius}
    px={{ base: '24px', md: '40px' }}
    py={{ base: '56px', md: '80px' }}
    textAlign="center"
  >
    <Text
      fontFamily={FONT}
      fontSize={{ base: '18px', md: '20px' }}
      fontWeight="700"
      color={COLOR.black}
    >
      {title}
    </Text>
    <Text fontFamily={FONT} fontSize="14px" color={COLOR.muted} mt="12px" lineHeight="1.8">
      {body}
    </Text>
  </Box>
);

export const PublicShare = () => {
  const { token } = useParams();
  const { treatments, isLoading, error } = usePublicShare(token);

  const isUnavailable = error instanceof ApiError && error.code === 'not_found';

  return (
    <Box display="flex" flexDirection="column" minHeight="100vh" bg={COLOR.white}>
      <Box
        as="header"
        bg={COLOR.white}
        borderBottom={`1px solid ${COLOR.border}`}
      >
        <Box
          maxW={PAGE_MAX_W}
          margin="0 auto"
          px={PAGE_PX}
          height={{ base: '60px', md: '72px' }}
          display="flex"
          alignItems="center"
        >
          <Text
            fontFamily={FONT}
            fontSize={{ base: '16px', md: '18px' }}
            fontWeight="700"
            letterSpacing="-0.01em"
            color={COLOR.black}
          >
            HairHistory
          </Text>
        </Box>
      </Box>

      <Box
        as="main"
        flex="1"
        width="100%"
        maxW={PAGE_MAX_W}
        margin="0 auto"
        px={PAGE_PX}
        py={{ base: '40px', md: '64px' }}
      >
        <PageHeader
          title="共有された施術履歴"
          description="リンクの持ち主が公開している施術の記録です。"
        />

        {isLoading && <SkeletonLoader />}

        {!isLoading && isUnavailable && (
          <EmptyState
            title="リンクが利用できません"
            body="この共有リンクは期限切れ・無効化済み、または存在しません。発行した方に新しいリンクを依頼してください。"
          />
        )}

        {!isLoading && error && !isUnavailable && (
          <EmptyState title="読み込みに失敗しました" body={error.message} />
        )}

        {!isLoading && !error && treatments.length === 0 && (
          <EmptyState title="記録がありません" body="共有された施術履歴はまだ登録されていません。" />
        )}

        {!isLoading && !error && treatments.length > 0 && (
          <Grid
            templateColumns={{ base: '1fr', md: 'repeat(2, 1fr)', lg: 'repeat(3, 1fr)' }}
            gap={{ base: '16px', md: '24px' }}
          >
            {treatments.map((treatment, index) => (
              <PublicTreatmentCard key={`${treatment.treatedOn}-${index}`} treatment={treatment} />
            ))}
          </Grid>
        )}
      </Box>
    </Box>
  );
};
