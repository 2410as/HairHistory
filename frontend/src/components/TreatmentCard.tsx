import { Box, Text } from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import { ServiceTags } from './ServiceTags';
import { CARD, COLOR, FONT } from '../design';
import { formatTreatedOn } from '../lib/format';
import type { Treatment } from '../types';

interface Props {
  treatment: Treatment;
}

export const TreatmentCard = ({ treatment }: Props) => {
  const navigate = useNavigate();

  return (
    <Box
      as="button"
      textAlign="left"
      width="100%"
      bg={CARD.bg}
      border={CARD.border}
      borderRadius={CARD.borderRadius}
      p={{ base: '20px', md: '24px' }}
      display="flex"
      flexDirection="column"
      minHeight="200px"
      transition="all 200ms ease"
      _hover={{
        borderColor: COLOR.borderStrong,
        boxShadow: '0 8px 24px rgba(0,0,0,0.08)',
        transform: 'translateY(-2px)',
      }}
      onClick={() => navigate(`/treatment/${treatment.id}`)}
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
        <Text
          fontFamily={FONT}
          fontSize="14px"
          color={COLOR.muted}
          lineHeight="1.8"
          mt="14px"
        >
          {treatment.memo}
        </Text>
      )}

      <Box
        mt="auto"
        pt="20px"
        width="100%"
        display="flex"
        justifyContent="space-between"
        alignItems="center"
      >
        <Text fontFamily={FONT} fontSize="15px" fontWeight="600" color={COLOR.black}>
          {typeof treatment.cost === 'number' ? `¥${treatment.cost.toLocaleString()}` : ''}
        </Text>
        <Text fontFamily={FONT} fontSize="13px" fontWeight="600" color={COLOR.muted}>
          詳細を見る →
        </Text>
      </Box>
    </Box>
  );
};
