import { chakra } from '@chakra-ui/react';
import type { ComponentProps } from 'react';
import { PRIMARY_BUTTON, SECONDARY_BUTTON } from '../design';

const StyledButton = chakra('button');

type Props = ComponentProps<typeof StyledButton> & {
  variant?: 'primary' | 'secondary';
};

export const ActionButton = ({ variant = 'primary', ...props }: Props) => (
  <StyledButton
    type="button"
    {...(variant === 'primary' ? PRIMARY_BUTTON : SECONDARY_BUTTON)}
    _disabled={{ opacity: 0.5, cursor: 'not-allowed', transform: 'none', boxShadow: 'none' }}
    {...props}
  />
);
