export class CreateProfileDto {
  userId: string;
  fileId?: string;
  fileUri?: string;
  fileThumbnailUri?: string;
  bankAccountName?: string;
  bankAccountHolder?: string;
  bankAccountNumber?: string;
}