import { Box, Input, Text } from '@chakra-ui/react';
import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { ActionButton } from '../components/ActionButton';
import { PageHeader } from '../components/PageHeader';
import { useAuth } from '../contexts/AuthContext';
import { CARD, COLOR, FONT } from '../design';
import { toErrorMessage } from '../lib/apiClient';

const FIELD_STYLES = {
  fontFamily: FONT,
  fontSize: '15px',
  height: 'auto',
  bg: COLOR.white,
  color: COLOR.black,
  border: `1px solid ${COLOR.border}`,
  borderRadius: '10px',
  px: '14px',
  py: '12px',
  transition: 'border-color 200ms ease, box-shadow 200ms ease',
  _placeholder: { color: COLOR.faint },
  _hover: { borderColor: COLOR.borderStrong },
  _focus: {
    borderColor: COLOR.black,
    boxShadow: '0 0 0 3px rgba(0,0,0,0.06)',
    outline: 'none',
  },
};

const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID ?? '';

interface LocationState {
  from?: string;
}

export const Login = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { devLogin } = useAuth();

  const [email, setEmail] = useState('anna@example.com');
  const [name, setName] = useState('Anna');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const redirectTo = (location.state as LocationState | null)?.from ?? '/dashboard';

  const handleDevLogin = async () => {
    if (!email.trim()) {
      setError('メールアドレスを入力してください');
      return;
    }
    setIsSubmitting(true);
    setError(null);
    try {
      await devLogin(email.trim(), name.trim() || email.trim());
      navigate(redirectTo, { replace: true });
    } catch (err) {
      setError(toErrorMessage(err));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Box>
      <PageHeader
        title="ログイン"
        description="施術履歴を記録・共有するにはログインが必要です。"
      />

      <Box maxW="480px">
        <Box
          bg={CARD.bg}
          border={CARD.border}
          borderRadius={CARD.borderRadius}
          p={{ base: '24px', md: '32px' }}
        >
          <ActionButton
            variant="secondary"
            width="100%"
            fontSize="15px"
            px="24px"
            py="14px"
            disabled={!GOOGLE_CLIENT_ID}
          >
            Google でログイン
          </ActionButton>
          {!GOOGLE_CLIENT_ID && (
            <Text fontFamily={FONT} fontSize="13px" color={COLOR.faint} mt="10px" lineHeight="1.7">
              クライアント ID 未設定のため無効です。VITE_GOOGLE_CLIENT_ID
              を設定すると有効になります。
            </Text>
          )}

          <Box
            display="flex"
            alignItems="center"
            gap="12px"
            my="28px"
            aria-hidden="true"
          >
            <Box flex="1" height="1px" bg={COLOR.border} />
            <Text fontFamily={FONT} fontSize="12px" color={COLOR.faint}>
              または
            </Text>
            <Box flex="1" height="1px" bg={COLOR.border} />
          </Box>

          <Text fontFamily={FONT} fontSize="14px" fontWeight="600" color={COLOR.black} mb="8px">
            開発用ログイン
          </Text>

          <Input
            aria-label="メールアドレス"
            placeholder="メールアドレス"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            width="100%"
            mb="12px"
            {...FIELD_STYLES}
          />
          <Input
            aria-label="表示名"
            placeholder="表示名"
            value={name}
            onChange={(event) => setName(event.target.value)}
            width="100%"
            {...FIELD_STYLES}
          />

          {error && (
            <Text fontFamily={FONT} fontSize="13px" color={COLOR.muted} mt="12px" lineHeight="1.7">
              {error}
            </Text>
          )}

          <ActionButton
            width="100%"
            fontSize="15px"
            px="24px"
            py="14px"
            mt="20px"
            disabled={isSubmitting}
            onClick={handleDevLogin}
          >
            {isSubmitting ? 'ログイン中…' : '開発用ログイン'}
          </ActionButton>
        </Box>
      </Box>
    </Box>
  );
};
