import Attachment, { CreateAttachment } from '@/api/models/Attachment';
import Message, { DefaultMessage } from '@/api/models/Message';
import RestResponse from '@/api/models/RestResponse';
import SessionApi from '@/api/SessionApi';
import { createClient } from '@/clients';
import Button from '@/components/common/Button';
import Card from '@/components/common/Card';
import IconButton from '@/components/common/IconButton';
import TextBox from '@/components/common/TextBox';
import {
  faEyeSlash,
  faFileAlt as faFile,
  faImage,
} from '@fortawesome/free-regular-svg-icons';
import {
  faEye,
  faFileAlt as faFileSolid,
  faImage as faImageSolid,
  faTimes,
} from '@fortawesome/free-solid-svg-icons';
import React, { useRef, useState } from 'react';
import styles from './styles.module.css';

const client = createClient();

type AttachmentMap = {
  image?: boolean;
  file?: boolean;
};

interface FileAttachment extends CreateAttachment {
  file?: File;
}

export interface SessionFormProps {
  api: SessionApi;
  sessionId: string;
}

const SessionForm: React.FC<SessionFormProps> = ({ api, sessionId }) => {
  const [valid, setValid] = useState(false);
  const [processing, setProcessing] = useState(false);
  const [textMessage, setTextMessage] = useState('');
  const [sensitive, setSensitive] = useState(false);
  const [imageSrc, setImageSrc] = useState('');
  const [attachment, setAttachment] = useState<AttachmentMap>({});

  const fileInputRef = useRef<HTMLInputElement>();
  const attachmentRef = useRef<FileAttachment>();

  const resetForm = () => {
    handleFileClear();
    setTextMessage('');
    setSensitive(false);
    setValid(false);
  };

  const handleValid = (body: string, file?: File) => {
    const nextValid = !!body || !!file;
    if (valid !== nextValid) setValid(nextValid);
  };

  const handleTextChange = (value: string) => {
    handleValid(value);
    setTextMessage(value);
  };

  const handleSend = async () => {
    if (!valid || processing) return;
    setProcessing(true);
    try {
      const message = DefaultMessage();
      message.sensitive = sensitive;
      message.body = textMessage;
      if (attachmentRef.current) {
        const attachment = await handleFileUpload(attachmentRef.current);
        message.attachment_id = attachment.id;
        message.attachment_name = attachment.name;
        message.attachment_type = attachment.type;
      }
      await api.sendMessage(sessionId, message as Message);
      resetForm();
    } catch (err) {
      console.error(err);
    } finally {
      setProcessing(false);
    }
  };

  const handleSensitive = () => setSensitive(!sensitive);

  const openFilePrompt = (image: boolean) => {
    const input = fileInputRef.current;
    input.value = '';
    input.accept = image ? 'image/*' : '*/*';
    input.click();
  };

  const handleFileChange = () => {
    const input = fileInputRef.current;
    const files = input.files;
    if (files.length <= 0) {
      setAttachment({});
      return;
    }
    const image = input.accept.startsWith('image/');
    const file = files[0];
    if (image) handleFileImage(file);
    attachmentRef.current = {
      file: file,
      name: file.name,
      type: image ? 'image' : 'file',
      mime: file.type,
    };
    setAttachment({ file: !image, image });
    handleValid(textMessage, file);
  };

  const handleFileClear = () => {
    fileInputRef.current.value = '';
    attachmentRef.current = undefined;
    setAttachment({});
    setImageSrc('');
  };

  const handleFileImage = (file: File) => {
    const reader = new FileReader();
    reader.onload = () => setImageSrc(reader.result as string);
    reader.onerror = (err) => console.error(err);
    reader.readAsDataURL(file);
  };

  const handleFileUpload = async (attachment: FileAttachment) => {
    if (!attachment.file) return;
    const formData = new FormData();
    formData.append('file', attachment.file, attachment.name);
    const { data } = await client.post('/attachments/upload', formData, {
      params: {
        type: attachment.type,
      },
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    const res = data as RestResponse<Attachment>;
    return res.data;
  };

  const handleImage = () => openFilePrompt(true);
  const handleFile = () => openFilePrompt(false);

  return (
    <Card className='card py-2 space-y-2'>
      <div className='px-2'>
        <TextBox
          value={textMessage}
          placeholder='Type your message here'
          onChange={handleTextChange}
        />
      </div>
      {imageSrc && (
        <div className='flex justify-center'>
          <img src={imageSrc} />
        </div>
      )}
      {attachmentRef.current && (
        <div className='flex justify-start items-center text-xs overflow-hidden px-1'>
          <IconButton icon={faTimes} color='red' onClick={handleFileClear} />
          <div className={styles.filename}>{attachmentRef.current.name}</div>
        </div>
      )}
      <div className='flex flex-row-reverse items-center px-1'>
        <div className='px-1'>
          <Button
            color='primary'
            className='rounded-full'
            onClick={handleSend}
            disabled={!valid || processing}
          >
            Send
          </Button>
        </div>
        <div>
          <IconButton
            icon={sensitive ? faEyeSlash : faEye}
            color={sensitive ? 'blue' : ''}
            onClick={handleSensitive}
          />
        </div>
        <div className='flex-grow'></div>
        <div>
          <IconButton
            icon={attachment.file ? faFileSolid : faFile}
            color={attachment.file ? 'blue' : ''}
            onClick={handleFile}
          />
        </div>
        <div>
          <IconButton
            icon={attachment.image ? faImageSolid : faImage}
            color={attachment.image ? 'blue' : ''}
            onClick={handleImage}
          />
        </div>
      </div>
      <input
        ref={fileInputRef}
        type='file'
        className='hidden'
        onChange={handleFileChange}
      />
    </Card>
  );
};

export default SessionForm;
