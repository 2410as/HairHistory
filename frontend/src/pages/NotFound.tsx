import { Box, Heading, Text } from '@chakra-ui/react'
import { useNavigate } from 'react-router-dom'
import { COLOR, FONT, PRIMARY_BUTTON } from '../design'

export function NotFound() {
  const navigate = useNavigate()

  return (
    <Box
      minHeight="100vh"
      display="flex"
      alignItems="center"
      justifyContent="center"
      bg={COLOR.white}
      px="24px"
    >
      <Box textAlign="center">
        <Heading
          as="h1"
          fontFamily={FONT}
          fontSize={{ base: '40px', md: '52px' }}
          fontWeight="700"
          letterSpacing="-0.01em"
          color={COLOR.black}
        >
          404
        </Heading>
        <Text
          fontFamily={FONT}
          fontSize={{ base: '14px', md: '16px' }}
          color={COLOR.muted}
          lineHeight="1.8"
          mt="14px"
        >
          お探しのページは見つかりませんでした。
        </Text>
        <Box display="flex" justifyContent="center" mt="32px">
          <Box
            as="button"
            {...PRIMARY_BUTTON}
            fontSize="15px"
            px="36px"
            py="15px"
            onClick={() => navigate('/')}
          >
            ホームに戻る
          </Box>
        </Box>
      </Box>
    </Box>
  )
}
