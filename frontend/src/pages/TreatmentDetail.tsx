import { Box, Input, Textarea, Text } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { ActionButton } from "../components/ActionButton";
import { PageHeader } from "../components/PageHeader";
import { SkeletonLoader } from "../components/SkeletonLoader";
import { CARD, COLOR, FONT, SECONDARY_BUTTON } from "../design";
import {
  useCreateTreatment,
  useDeleteTreatment,
  useTreatment,
  useUpdateTreatment,
} from "../hooks/useTreatments";
import { toErrorMessage } from "../lib/apiClient";
import type { TreatmentInput } from "../types";

const FIELD_STYLES = {
  fontFamily: FONT,
  fontSize: "15px",
  height: "auto",
  bg: COLOR.white,
  color: COLOR.black,
  border: `1px solid ${COLOR.border}`,
  borderRadius: "10px",
  px: "14px",
  py: "12px",
  transition: "border-color 200ms ease, box-shadow 200ms ease",
  _placeholder: { color: COLOR.faint },
  _hover: { borderColor: COLOR.borderStrong },
  _focus: {
    borderColor: COLOR.black,
    boxShadow: "0 0 0 3px rgba(0,0,0,0.06)",
    outline: "none",
  },
};

const SERVICE_OPTIONS = ["カット", "カラー", "パーマ", "トリートメント", "ヘッドスパ"];

const MAX_SERVICES = 10;
const MAX_SERVICE_LENGTH = 50;

interface FieldProps {
  label: string;
  required?: boolean;
  error?: string;
  children: React.ReactNode;
}

const Field = ({ label, required, error, children }: FieldProps) => (
  <Box mb="24px">
    <Text
      as="label"
      display="block"
      fontFamily={FONT}
      fontSize="14px"
      fontWeight="600"
      color={COLOR.black}
      mb="8px"
    >
      {label}
      {required && (
        <Text as="span" color={COLOR.faint} fontWeight="500" ml="6px" fontSize="12px">
          必須
        </Text>
      )}
    </Text>
    {children}
    {error && (
      <Text fontFamily={FONT} fontSize="13px" color={COLOR.muted} mt="8px">
        {error}
      </Text>
    )}
  </Box>
);

interface ServiceSelectorProps {
  selected: string[];
  onChange: (services: string[]) => void;
}

const ServiceSelector = ({ selected, onChange }: ServiceSelectorProps) => {
  const [draft, setDraft] = useState("");
  const [draftError, setDraftError] = useState("");

  const isFull = selected.length >= MAX_SERVICES;
  const customServices = selected.filter((service) => !SERVICE_OPTIONS.includes(service));

  const toggleOption = (option: string) => {
    setDraftError("");
    if (selected.includes(option)) {
      onChange(selected.filter((service) => service !== option));
      return;
    }
    if (isFull) {
      setDraftError(`施術項目は最大 ${MAX_SERVICES} 件までです`);
      return;
    }
    onChange([...selected, option]);
  };

  const addDraft = () => {
    const trimmed = draft.trim();
    if (!trimmed) {
      setDraftError("追加する施術項目を入力してください");
      return;
    }
    if (trimmed.length > MAX_SERVICE_LENGTH) {
      setDraftError(`1 件あたり ${MAX_SERVICE_LENGTH} 文字までです`);
      return;
    }
    if (selected.includes(trimmed)) {
      setDraftError("すでに追加されています");
      return;
    }
    if (isFull) {
      setDraftError(`施術項目は最大 ${MAX_SERVICES} 件までです`);
      return;
    }
    onChange([...selected, trimmed]);
    setDraft("");
    setDraftError("");
  };

  const removeService = (service: string) => {
    setDraftError("");
    onChange(selected.filter((item) => item !== service));
  };

  return (
    <Box>
      <Box display="flex" flexWrap="wrap" gap="8px">
        {SERVICE_OPTIONS.map((option) => {
          const isSelected = selected.includes(option);
          return (
            <Box
              as="button"
              key={option}
              role="checkbox"
              aria-checked={isSelected}
              aria-label={option}
              onClick={() => toggleOption(option)}
              fontFamily={FONT}
              fontSize="14px"
              fontWeight="600"
              borderRadius="999px"
              px="16px"
              py="9px"
              cursor="pointer"
              transition="all 200ms ease"
              bg={isSelected ? COLOR.black : COLOR.white}
              color={isSelected ? COLOR.white : COLOR.black}
              border={`1px solid ${isSelected ? COLOR.black : COLOR.border}`}
              _hover={{
                borderColor: isSelected ? COLOR.black : COLOR.borderStrong,
                bg: isSelected ? "#333333" : COLOR.surface,
              }}
            >
              {option}
            </Box>
          );
        })}
      </Box>

      {customServices.length > 0 && (
        <Box display="flex" flexWrap="wrap" gap="8px" mt="12px">
          {customServices.map((service) => (
            <Box
              key={service}
              bg={COLOR.surface}
              borderRadius="999px"
              pl="14px"
              pr="8px"
              py="7px"
              display="flex"
              alignItems="center"
              gap="8px"
            >
              <Text fontFamily={FONT} fontSize="13px" fontWeight="600" color={COLOR.black}>
                {service}
              </Text>
              <Box
                as="button"
                aria-label={`${service} を外す`}
                onClick={() => removeService(service)}
                fontFamily={FONT}
                fontSize="13px"
                lineHeight="1"
                color={COLOR.muted}
                bg="transparent"
                border="none"
                cursor="pointer"
                px="4px"
                _hover={{ color: COLOR.black }}
              >
                ×
              </Box>
            </Box>
          ))}
        </Box>
      )}

      <Box display="flex" gap="8px" mt="12px">
        <Input
          aria-label="その他の施術項目"
          value={draft}
          maxLength={MAX_SERVICE_LENGTH}
          onChange={(e) => {
            setDraft(e.target.value);
            setDraftError("");
          }}
          onKeyDown={(e) => {
            if (e.key !== "Enter") return;
            e.preventDefault();
            addDraft();
          }}
          placeholder="例）縮毛矯正、ブリーチ"
          flex="1"
          {...FIELD_STYLES}
        />
        <Box
          as="button"
          {...SECONDARY_BUTTON}
          fontSize="14px"
          px="20px"
          py="12px"
          flexShrink={0}
          onClick={addDraft}
        >
          追加
        </Box>
      </Box>

      <Text fontFamily={FONT} fontSize="12px" color={COLOR.faint} mt="8px" lineHeight="1.7">
        {`選択中 ${selected.length} / ${MAX_SERVICES} 件。1 件あたり ${MAX_SERVICE_LENGTH} 文字まで。`}
      </Text>

      {draftError && (
        <Text fontFamily={FONT} fontSize="13px" color={COLOR.muted} mt="6px">
          {draftError}
        </Text>
      )}
    </Box>
  );
};

interface FormState {
  services: string[];
  treatedOn: string;
  salonName: string;
  memo: string;
  cost: string;
}

const EMPTY_FORM: FormState = {
  services: [],
  treatedOn: "",
  salonName: "",
  memo: "",
  cost: "",
};

export const TreatmentDetail = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const isNewTreatment = id === "new";

  const { data: treatment, isLoading, error: loadError } = useTreatment(
    isNewTreatment ? undefined : id
  );
  const createMutation = useCreateTreatment();
  const updateMutation = useUpdateTreatment();
  const deleteMutation = useDeleteTreatment();

  const [formData, setFormData] = useState<FormState>(EMPTY_FORM);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [isConfirmingDelete, setIsConfirmingDelete] = useState(false);

  useEffect(() => {
    if (!treatment) return;
    setFormData({
      services: [...treatment.services],
      treatedOn: treatment.treatedOn,
      salonName: treatment.salonName ?? "",
      memo: treatment.memo ?? "",
      cost: typeof treatment.cost === "number" ? String(treatment.cost) : "",
    });
  }, [treatment]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: "" }));
    }
  };

  const setServices = (services: string[]) => {
    setFormData((prev) => ({ ...prev, services }));
    if (errors.services) {
      setErrors((prev) => ({ ...prev, services: "" }));
    }
  };

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};
    const services = formData.services.map((service) => service.trim());
    if (services.length === 0) {
      newErrors.services = "施術項目を 1 件以上選択してください";
    } else if (services.length > MAX_SERVICES) {
      newErrors.services = `施術項目は最大 ${MAX_SERVICES} 件までです`;
    } else if (services.some((service) => !service)) {
      newErrors.services = "空の施術項目は保存できません";
    } else if (services.some((service) => service.length > MAX_SERVICE_LENGTH)) {
      newErrors.services = `施術項目は 1 件あたり ${MAX_SERVICE_LENGTH} 文字までです`;
    }
    if (!formData.treatedOn) {
      newErrors.treatedOn = "日付を入力してください";
    }
    if (formData.cost && !/^\d+$/.test(formData.cost.trim())) {
      newErrors.cost = "料金は 0 以上の整数で入力してください";
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const buildInput = (): TreatmentInput => {
    const input: TreatmentInput = {
      treatedOn: formData.treatedOn,
      services: formData.services.map((service) => service.trim()),
    };
    if (formData.salonName.trim()) input.salonName = formData.salonName.trim();
    if (formData.memo.trim()) input.memo = formData.memo.trim();
    if (formData.cost.trim()) input.cost = Number(formData.cost.trim());
    return input;
  };

  const handleSave = async () => {
    setSubmitError(null);
    if (!validate()) return;

    const input = buildInput();
    const saved = isNewTreatment
      ? await createMutation.run(input)
      : await updateMutation.run(id as string, input);

    if (saved) {
      navigate("/dashboard");
      return;
    }

    const failure = isNewTreatment ? createMutation.error : updateMutation.error;
    setSubmitError(toErrorMessage(failure));
  };

  const handleDelete = async () => {
    setSubmitError(null);
    const deleted = await deleteMutation.run(id as string);
    if (deleted) {
      navigate("/dashboard");
      return;
    }
    setIsConfirmingDelete(false);
    setSubmitError(toErrorMessage(deleteMutation.error));
  };

  const isSubmitting =
    createMutation.isSubmitting || updateMutation.isSubmitting || deleteMutation.isSubmitting;

  if (!isNewTreatment && isLoading) {
    return (
      <Box>
        <PageHeader title="施術の詳細" description="記録した施術の内容を確認・編集できます。" />
        <SkeletonLoader />
      </Box>
    );
  }

  if (!isNewTreatment && loadError) {
    return (
      <Box>
        <PageHeader title="施術の詳細" description="記録した施術の内容を確認・編集できます。" />
        <Box
          bg={CARD.bg}
          border={CARD.border}
          borderRadius={CARD.borderRadius}
          px={{ base: "24px", md: "40px" }}
          py={{ base: "48px", md: "64px" }}
          textAlign="center"
        >
          <Text fontFamily={FONT} fontSize="16px" fontWeight="600" color={COLOR.black}>
            読み込みに失敗しました
          </Text>
          <Text fontFamily={FONT} fontSize="14px" color={COLOR.muted} mt="10px" lineHeight="1.8">
            {toErrorMessage(loadError)}
          </Text>
          <Box display="flex" justifyContent="center" mt="28px">
            <Box
              as="button"
              {...SECONDARY_BUTTON}
              fontSize="15px"
              px="32px"
              py="14px"
              onClick={() => navigate("/dashboard")}
            >
              一覧に戻る
            </Box>
          </Box>
        </Box>
      </Box>
    );
  }

  return (
    <Box>
      <PageHeader
        title={isNewTreatment ? "施術を追加" : "施術の詳細"}
        description={
          isNewTreatment
            ? "サロンでの施術内容を記録しておくと、次回の相談がスムーズになります。"
            : "記録した施術の内容を確認・編集できます。"
        }
      />

      <Box maxW="640px">
        <Box
          bg={CARD.bg}
          border={CARD.border}
          borderRadius={CARD.borderRadius}
          p={{ base: "24px", md: "32px" }}
        >
          <Field label="施術項目" required error={errors.services}>
            <ServiceSelector selected={formData.services} onChange={setServices} />
          </Field>

          <Field label="日付" required error={errors.treatedOn}>
            <Input
              type="date"
              name="treatedOn"
              aria-label="日付"
              value={formData.treatedOn}
              onChange={handleChange}
              width="100%"
              {...FIELD_STYLES}
            />
          </Field>

          <Field label="サロン名">
            <Input
              name="salonName"
              aria-label="サロン名"
              value={formData.salonName}
              onChange={handleChange}
              placeholder="例）Hair Salon ABC"
              width="100%"
              {...FIELD_STYLES}
            />
          </Field>

          <Field label="料金" error={errors.cost}>
            <Input
              name="cost"
              aria-label="料金"
              inputMode="numeric"
              value={formData.cost}
              onChange={handleChange}
              placeholder="例）8000"
              width="100%"
              {...FIELD_STYLES}
            />
          </Field>

          <Field label="メモ">
            <Textarea
              name="memo"
              aria-label="メモ"
              value={formData.memo}
              onChange={handleChange}
              placeholder="薬剤・仕上がり・次回の希望など"
              width="100%"
              minHeight="128px"
              resize="none"
              lineHeight="1.8"
              {...FIELD_STYLES}
            />
          </Field>

          {submitError && (
            <Text fontFamily={FONT} fontSize="13px" color={COLOR.muted} lineHeight="1.7">
              {submitError}
            </Text>
          )}
        </Box>

        <Box
          display="flex"
          flexDirection={{ base: "column-reverse", sm: "row" }}
          justifyContent="flex-end"
          gap="12px"
          mt="28px"
        >
          {!isNewTreatment && (
            <ActionButton
              variant="secondary"
              fontSize="15px"
              px="32px"
              py="14px"
              marginRight={{ sm: "auto" }}
              disabled={isSubmitting}
              onClick={() => setIsConfirmingDelete(true)}
            >
              削除する
            </ActionButton>
          )}
          <ActionButton
            variant="secondary"
            fontSize="15px"
            px="32px"
            py="14px"
            disabled={isSubmitting}
            onClick={() => navigate("/dashboard")}
          >
            キャンセル
          </ActionButton>
          <ActionButton
            fontSize="15px"
            px="32px"
            py="14px"
            disabled={isSubmitting}
            onClick={handleSave}
          >
            {isSubmitting ? "保存中…" : "保存する"}
          </ActionButton>
        </Box>

        {isConfirmingDelete && (
          <Box
            bg={CARD.bg}
            border={CARD.border}
            borderRadius={CARD.borderRadius}
            p={{ base: "20px", md: "24px" }}
            mt="20px"
            role="alertdialog"
            aria-label="削除の確認"
          >
            <Text fontFamily={FONT} fontSize="15px" fontWeight="700" color={COLOR.black}>
              この施術を削除しますか？
            </Text>
            <Text fontFamily={FONT} fontSize="14px" color={COLOR.muted} mt="8px" lineHeight="1.8">
              削除すると元に戻せません。
            </Text>
            <Box display="flex" gap="12px" mt="20px" flexWrap="wrap">
              <ActionButton
                fontSize="14px"
                px="26px"
                py="12px"
                disabled={isSubmitting}
                onClick={handleDelete}
              >
                削除を確定
              </ActionButton>
              <ActionButton
                variant="secondary"
                fontSize="14px"
                px="26px"
                py="12px"
                disabled={isSubmitting}
                onClick={() => setIsConfirmingDelete(false)}
              >
                やめる
              </ActionButton>
            </Box>
          </Box>
        )}
      </Box>
    </Box>
  );
};
