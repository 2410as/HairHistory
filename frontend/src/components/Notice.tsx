import { Box, Text } from '@chakra-ui/react';
import type { ReactNode } from 'react';
import { COLOR, FONT, NOTICE } from '../design';

export const Notice = ({ children }: { children: ReactNode }) => (
  <Box
    role="alert"
    display="flex"
    alignItems="flex-start"
    gap="10px"
    bg={NOTICE.bg}
    color={NOTICE.color}
    borderRadius={NOTICE.borderRadius}
    px={NOTICE.px}
    py={NOTICE.py}
  >
    <Box
      aria-hidden="true"
      flexShrink={0}
      width="18px"
      height="18px"
      mt="2px"
      borderRadius="999px"
      border={`1.5px solid ${COLOR.white}`}
      display="flex"
      alignItems="center"
      justifyContent="center"
      fontFamily={FONT}
      fontSize="12px"
      fontWeight="700"
      lineHeight="1"
    >
      !
    </Box>
    <Text fontFamily={FONT} fontSize={NOTICE.fontSize} fontWeight="600" lineHeight="1.7">
      {children}
    </Text>
  </Box>
);
