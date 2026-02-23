# Backend API - Promoter, Material, AI Endpoints

## ADDED Requirements

### Requirement: Promoter List API
The system SHALL provide a promoter list API endpoint that supports pagination and filtering.

#### Scenario: Query promoter list
- Given a brand admin is logged in
- When requesting GET /promoter/list?page=1&pageSize=20
- Then return promoter list with id, name, phone, status, level, totalOrders, totalRewards, conversionRate

#### Scenario: Filter promoters by status
- Given a brand admin is logged in
- When requesting GET /promoter/list?status=active
- Then return only promoters with status=active

### Requirement: Promoter Detail API
The system SHALL provide a promoter detail API endpoint.

#### Scenario: Query promoter detail
- Given a brand admin is logged in
- When requesting GET /promoter/detail/123
- Then return full promoter info including basic info, performance stats, related campaigns, and links

### Requirement: Promoter Link Generation API
The system SHALL provide a promoter link generation API endpoint.

#### Scenario: Generate promoter link
- Given a brand admin is logged in
- When requesting POST /promoter/generate-link with body {promoterId: 123, campaignId: 9}
- Then return generated promoter link and QR code

### Requirement: Promoter Rewards API
The system SHALL provide a promoter rewards API endpoint.

#### Scenario: Query promoter rewards
- Given a brand admin is logged in
- When requesting GET /promoter/rewards/123?page=1&pageSize=20
- Then return reward list with id, type, status, amount, description, createdAt

### Requirement: Material Upload API
The system SHALL provide a material upload API endpoint.

#### Scenario: Upload image material
- Given a brand admin is logged in
- When requesting POST /material/upload with multipart/form-data (file, type=image)
- Then save material and return id, url, name

### Requirement: Material List API
The system SHALL provide a material list API endpoint.

#### Scenario: Query material list
- Given a brand admin is logged in
- When requesting GET /material/list?page=1&pageSize=20
- Then return material list with id, name, description, type, url, createdAt

### Requirement: Material Delete API
The system SHALL provide a material delete API endpoint.

#### Scenario: Delete material
- Given a brand admin is logged in
- When requesting DELETE /material/delete/123
- Then delete material file and metadata, return success

### Requirement: AI Copywriting Generation API
The system SHALL provide an AI copywriting generation API endpoint.

#### Scenario: Generate marketing copy
- Given a brand admin is logged in
- When requesting POST /ai/generate-copywriting with body {topic: sale, style: professional, length: medium}
- Then return generated copywriting content

#### Scenario: Fallback when AI service unavailable
- Given AI service is unavailable
- When requesting POST /ai/generate-copywriting
- Then return template-based copywriting without error
