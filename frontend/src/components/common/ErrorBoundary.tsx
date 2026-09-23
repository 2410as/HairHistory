import React, { type ReactNode, type ErrorInfo } from 'react';
import { Box, Heading, Text } from '@chakra-ui/react';
import { COLOR, FONT, PRIMARY_BUTTON } from '../../design';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
  errorInfo?: ErrorInfo;
}

export class ErrorBoundary extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): Partial<State> {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    // Update state to capture errorInfo for logging/debugging
    this.setState({ errorInfo });
    // Log error details (can be integrated with error tracking service)
    console.error('Error caught by boundary:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback error={this.state.error} />;
    }

    return this.props.children;
  }
}

interface ErrorFallbackProps {
  error?: Error;
}

const ErrorFallback = ({ error }: ErrorFallbackProps) => {
  const handleReload = () => {
    window.location.reload();
  };

  return (
    <Box
      bg={COLOR.white}
      minHeight="100vh"
      display="flex"
      alignItems="center"
      justifyContent="center"
      py="40px"
      px="24px"
    >
      <Box textAlign="center" maxW="480px">
        <Heading
          as="h1"
          fontFamily={FONT}
          fontSize={{ base: '24px', md: '30px' }}
          fontWeight="700"
          letterSpacing="-0.01em"
          color={COLOR.black}
        >
          問題が発生しました
        </Heading>
        <Text
          fontFamily={FONT}
          fontSize={{ base: '14px', md: '15px' }}
          color={COLOR.muted}
          lineHeight="1.8"
          mt="14px"
        >
          {error?.message || '予期しないエラーが発生しました。'}
        </Text>
        <Box display="flex" justifyContent="center" mt="32px">
          <Box
            as="button"
            {...PRIMARY_BUTTON}
            fontSize="15px"
            px="36px"
            py="15px"
            onClick={handleReload}
            aria-label="ページを再読み込みする"
          >
            再読み込み
          </Box>
        </Box>
      </Box>
    </Box>
  );
};
