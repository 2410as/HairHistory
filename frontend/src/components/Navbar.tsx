import { Box, Flex, HStack, Text } from '@chakra-ui/react';
import { Link as RouterLink, useLocation, useNavigate } from 'react-router-dom';
import { ActionButton } from './ActionButton';
import { useAuth } from '../contexts/AuthContext';
import { COLOR, FONT, PAGE_MAX_W, PAGE_PX, SECONDARY_BUTTON } from '../design';

const NAV_ITEMS = [
  { label: 'ホーム', to: '/' },
  { label: '履歴', to: '/dashboard' },
  { label: '共有', to: '/share' },
];

export const Navbar = () => {
  const { pathname } = useLocation();
  const navigate = useNavigate();
  const { user, isLoading, logout } = useAuth();

  const handleLogout = async () => {
    await logout();
    navigate('/login', { replace: true });
  };

  return (
    <Box
      as="header"
      bg={COLOR.white}
      borderBottom={`1px solid ${COLOR.border}`}
      position="sticky"
      top="0"
      zIndex={20}
    >
      <Flex
        maxW={PAGE_MAX_W}
        margin="0 auto"
        px={PAGE_PX}
        height={{ base: 'auto', md: '72px' }}
        py={{ base: '12px', md: '0' }}
        flexWrap="wrap"
        rowGap="10px"
        justify="space-between"
        align="center"
        gap="16px"
      >
        <RouterLink to="/">
          <Text
            fontFamily={FONT}
            fontSize={{ base: '16px', md: '18px' }}
            fontWeight="700"
            letterSpacing="-0.01em"
            color={COLOR.black}
          >
            HairHistory
          </Text>
        </RouterLink>

        <HStack gap={{ base: '18px', md: '28px' }}>
          {NAV_ITEMS.map((item) => {
            const isActive =
              item.to === '/' ? pathname === '/' : pathname.startsWith(item.to);

            return (
              <RouterLink key={item.to} to={item.to}>
                <Text
                  fontFamily={FONT}
                  fontSize={{ base: '13px', md: '14px' }}
                  fontWeight={isActive ? '600' : '500'}
                  color={isActive ? COLOR.black : COLOR.muted}
                  transition="color 200ms ease"
                  _hover={{ color: COLOR.black }}
                  whiteSpace="nowrap"
                >
                  {item.label}
                </Text>
              </RouterLink>
            );
          })}
        </HStack>

        {!isLoading && (
          <HStack
            gap={{ base: '12px', md: '16px' }}
            order={{ base: 3, md: 0 }}
            width={{ base: '100%', md: 'auto' }}
            justifyContent={{ base: 'flex-end', md: 'flex-start' }}
          >
            {user ? (
              <>
                <Text
                  fontFamily={FONT}
                  fontSize={{ base: '13px', md: '14px' }}
                  fontWeight="600"
                  color={COLOR.black}
                  maxW={{ base: '150px', md: '160px' }}
                  overflow="hidden"
                  textOverflow="ellipsis"
                  whiteSpace="nowrap"
                  data-testid="current-user"
                >
                  {user.name}
                </Text>
                <ActionButton
                  variant="secondary"
                  fontSize={{ base: '12px', md: '13px' }}
                  px={{ base: '16px', md: '18px' }}
                  py={{ base: '8px', md: '9px' }}
                  whiteSpace="nowrap"
                  onClick={handleLogout}
                >
                  ログアウト
                </ActionButton>
              </>
            ) : (
              <RouterLink to="/login">
                <Box
                  as="span"
                  display="inline-block"
                  {...SECONDARY_BUTTON}
                  fontSize={{ base: '12px', md: '13px' }}
                  px={{ base: '16px', md: '18px' }}
                  py={{ base: '8px', md: '9px' }}
                  whiteSpace="nowrap"
                >
                  ログイン
                </Box>
              </RouterLink>
            )}
          </HStack>
        )}
      </Flex>
    </Box>
  );
};
