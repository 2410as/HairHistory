import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { ChakraProvider, createSystem, defaultConfig } from '@chakra-ui/react'
import { theme } from './theme'
import { ErrorBoundary } from './components/common/ErrorBoundary'
import Layout from './components/Layout'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Home } from './pages/Home'
import { Login } from './pages/Login'
import { Dashboard } from './pages/Dashboard'
import { TreatmentDetail } from './pages/TreatmentDetail'
import { TreatmentList } from './pages/TreatmentList'
import { ShareLink } from './pages/ShareLink'
import { PublicShare } from './pages/PublicShare'
import { NotFound } from './pages/NotFound'

const system = createSystem(defaultConfig, theme)

function App() {
  return (
    <ChakraProvider value={system}>
      <ErrorBoundary>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/share/:token" element={<PublicShare />} />
            <Route element={<Layout />}>
              <Route path="/login" element={<Login />} />
              <Route element={<ProtectedRoute />}>
                <Route path="/dashboard" element={<Dashboard />} />
                <Route path="/treatments" element={<TreatmentList />} />
                <Route path="/treatment/:id" element={<TreatmentDetail />} />
                <Route path="/share" element={<ShareLink />} />
              </Route>
            </Route>
            <Route path="*" element={<NotFound />} />
          </Routes>
        </BrowserRouter>
      </ErrorBoundary>
    </ChakraProvider>
  )
}

export default App
