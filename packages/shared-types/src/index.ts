// Core Domain Types for Social Platform

export type Role =
  | 'ROLE_ANONYMOUS'
  | 'ROLE_USER_STANDARD'
  | 'ROLE_USER_MINOR'
  | 'ROLE_CREATOR_PRO'
  | 'ROLE_PAGE_OWNER'
  | 'ROLE_PAGE_EDITOR'
  | 'ROLE_GROUP_ADMIN'
  | 'ROLE_GROUP_MOD'
  | 'ROLE_MARKET_SELLER'
  | 'ROLE_ADVERTISER'
  | 'ROLE_SUPPORT_AGENT'
  | 'ROLE_MODERATOR_SAFETY'
  | 'ROLE_PLATFORM_ADMIN';

export type UserStatus = 'ACTIVE' | 'SUSPENDED' | 'DEACTIVATED' | 'PENDING_VERIFICATION';

export interface User {
  id: string;
  email: string;
  phone?: string;
  status: UserStatus;
  emailVerified: boolean;
  phoneVerified: boolean;
  mfaEnabled: boolean;
  role: Role;
  createdAt: string;
  updatedAt: string;
}

export interface Profile {
  userId: string;
  username: string;
  firstName: string;
  lastName: string;
  bio?: string;
  avatarUrl?: string;
  coverUrl?: string;
  dateOfBirth: string;
  gender: 'MALE' | 'FEMALE' | 'CUSTOM';
  privacyLevel: 'PUBLIC' | 'FRIENDS' | 'ONLY_ME';
  isProfessionalMode: boolean;
}

export type RelationshipType = 'FRIEND' | 'FOLLOW' | 'BLOCK' | 'RESTRICT';
export type RelationshipStatus = 'PENDING' | 'ACCEPTED' | 'DECLINED';

export interface Relationship {
  id: string;
  userId: string;
  targetUserId: string;
  relationshipType: RelationshipType;
  status: RelationshipStatus;
  createdAt: string;
}

export type PrivacyLevel = 'PUBLIC' | 'FRIENDS' | 'ONLY_ME' | 'CUSTOM';
export type PostType = 'STANDARD' | 'REEL' | 'STORY' | 'LIVE_ARCHIVE';

export interface PostAttachment {
  id: string;
  postId: string;
  mediaUrl: string;
  mediaType: 'IMAGE' | 'VIDEO' | 'AUDIO' | 'DOCUMENT';
  width?: number;
  height?: number;
  durationSeconds?: number;
  sortOrder: number;
}

export interface Post {
  id: string;
  authorId: string;
  authorProfile?: Profile;
  groupId?: string;
  pageId?: string;
  content: string;
  privacy: PrivacyLevel;
  postType: PostType;
  attachments?: PostAttachment[];
  reactionCount: number;
  commentCount: number;
  shareCount: number;
  isPinned: boolean;
  createdAt: string;
  updatedAt: string;
}

export type ReactionType = 'LIKE' | 'LOVE' | 'CARE' | 'HAHA' | 'WOW' | 'SAD' | 'ANGRY';

export interface Reaction {
  id: string;
  targetId: string;
  targetType: 'POST' | 'COMMENT';
  userId: string;
  reactionType: ReactionType;
  createdAt: string;
}

export interface Comment {
  id: string;
  postId: string;
  authorId: string;
  authorProfile?: Profile;
  parentCommentId?: string;
  content: string;
  reactionCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  tokenType: 'Bearer';
}

export interface LoginResponse {
  user: User;
  profile: Profile;
  tokens: AuthTokens;
}
