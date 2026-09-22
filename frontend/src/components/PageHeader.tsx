import { Box, Heading, Text } from '@chakra-ui/react';
import { COLOR, FONT } from '../design';

interface Props {
  title: string;
  description?: string;
  action?: React.ReactNode;
}

export const PageHeader = ({ title, description, action }: Props) => (
  <Box
    display="flex"
    flexDirection={{ base: 'column', md: 'row' }}
    alignItems={{ base: 'flex-start', md: 'flex-end' }}
    justifyContent="space-between"
    gap={{ base: '20px', md: '32px' }}
    pb={{ base: '24px', md: '28px' }}
    borderBottom={`1px solid ${COLOR.border}`}
    mb={{ base: '32px', md: '44px' }}
  >
    <Box>
      <Heading
        as="h1"
        fontFamily={FONT}
        fontSize={{ base: '28px', md: '38px' }}
        fontWeight="700"
        lineHeight="1.3"
        letterSpacing="-0.01em"
        color={COLOR.black}
      >
        {title}
      </Heading>
      {description && (
        <Text
          fontFamily={FONT}
          fontSize={{ base: '14px', md: '15px' }}
          color={COLOR.muted}
          lineHeight="1.8"
          mt="10px"
        >
          {description}
        </Text>
      )}
    </Box>
    {action}
  </Box>
);
