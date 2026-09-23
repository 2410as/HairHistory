import { Box, Text } from '@chakra-ui/react';
import { COLOR, FONT } from '../design';

interface Props {
  services: string[];
}

export const ServiceTags = ({ services }: Props) => (
  <Box display="flex" flexWrap="wrap" gap="8px">
    {services.map((service) => (
      <Box key={service} bg={COLOR.surface} borderRadius="999px" px="12px" py="5px">
        <Text fontFamily={FONT} fontSize="12px" fontWeight="600" color={COLOR.black}>
          {service}
        </Text>
      </Box>
    ))}
  </Box>
);
