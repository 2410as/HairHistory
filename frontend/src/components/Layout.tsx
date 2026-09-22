import { Outlet } from 'react-router-dom'
import { Box } from '@chakra-ui/react'
import { Navbar } from './Navbar'
import { AuthProvider } from '../contexts/AuthContext'
import { COLOR, PAGE_MAX_W, PAGE_PX } from '../design'

export default function Layout() {
  return (
    <AuthProvider>
      <Box display="flex" flexDirection="column" minHeight="100vh" bg={COLOR.white}>
        <Navbar />
        <Box
          as="main"
          flex="1"
          width="100%"
          maxW={PAGE_MAX_W}
          margin="0 auto"
          px={PAGE_PX}
          py={{ base: '40px', md: '64px' }}
        >
          <Outlet />
        </Box>
      </Box>
    </AuthProvider>
  )
}
