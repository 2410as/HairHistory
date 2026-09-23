import { Box, Grid } from '@chakra-ui/react';
import { CARD, COLOR } from '../design';

const Bar = ({ width, height = '14px' }: { width: string; height?: string }) => (
  <Box className="skeleton-pulse" bg={COLOR.surface} borderRadius="6px" width={width} height={height} />
);

export const SkeletonLoader = () => (
  <Grid
    templateColumns={{ base: '1fr', md: 'repeat(2, 1fr)', lg: 'repeat(3, 1fr)' }}
    gap={{ base: '16px', md: '24px' }}
  >
    {[1, 2, 3].map((index) => (
      <Box
        key={index}
        bg={CARD.bg}
        border={CARD.border}
        borderRadius={CARD.borderRadius}
        p={{ base: '20px', md: '24px' }}
        minHeight="200px"
        display="flex"
        flexDirection="column"
        gap="14px"
      >
        <Bar width="72px" height="24px" />
        <Bar width="60%" height="18px" />
        <Bar width="40%" />
        <Bar width="85%" />
        <Box mt="auto">
          <Bar width="50%" />
        </Box>
      </Box>
    ))}
  </Grid>
);
