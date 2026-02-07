# campaign-management Specification

## Purpose
TBD - created by archiving change add-h5-campaign-page-designer. Update Purpose after archive.
## Requirements
### Requirement: H5 Campaign Page Designer
The H5 application SHALL provide a campaign page designer that allows brand administrators to design campaign landing pages with the same functionality as the Admin backend.

#### Scenario: Access page designer
- **WHEN** a brand administrator edits a campaign
- **THEN** they SHALL be able to access the page designer
- **AND** see a mobile-optimized design interface

#### Scenario: Add page components
- **WHEN** a brand administrator wants to add a component
- **THEN** they SHALL be able to select from a component library
- **AND** the component SHALL be added to the design canvas
- **AND** receive visual feedback of the addition

#### Scenario: Edit component content
- **WHEN** a brand administrator clicks on a component
- **THEN** a bottom sheet editor SHALL appear
- **AND** they SHALL be able to edit component properties
- **AND** changes SHALL be reflected immediately

#### Scenario: Reorder components
- **WHEN** a brand administrator wants to change component order
- **THEN** they SHALL be able to use up/down arrows
- **AND** the component SHALL move to the new position
- **AND** other components SHALL adjust accordingly

#### Scenario: Delete components
- **WHEN** a brand administrator wants to remove a component
- **THEN** they SHALL be able to delete it
- **AND** receive confirmation before deletion
- **AND** the component SHALL be removed from the canvas

#### Scenario: Preview page design
- **WHEN** a brand administrator clicks preview
- **THEN** they SHALL see a full-screen preview
- **AND** the preview SHALL show the actual mobile appearance
- **AND** all components SHALL be rendered correctly

#### Scenario: Save page configuration
- **WHEN** a brand administrator saves the page design
- **THEN** the configuration SHALL be persisted to the backend
- **AND** they SHALL receive confirmation of successful save
- **AND** the configuration SHALL be associated with the campaign

#### Scenario: Load existing configuration
- **WHEN** a brand administrator opens an existing campaign
- **THEN** the page designer SHALL load the saved configuration
- **AND** all components SHALL be restored correctly
- **AND** component order SHALL be preserved

### Requirement: Supported Page Components
The page designer SHALL support the following component types with mobile-optimized editing.

#### Scenario: Activity poster component
- **WHEN** adding an activity poster
- **THEN** the administrator SHALL be able to upload or select an image
- **AND** the image SHALL be displayed with proper aspect ratio
- **AND** support image URL input

#### Scenario: Activity title component
- **WHEN** adding an activity title
- **THEN** the administrator SHALL be able to enter title text
- **AND** optionally enter subtitle text
- **AND** use AI color adjustment feature

#### Scenario: Activity time component
- **WHEN** adding activity time
- **THEN** the administrator SHALL be able to select date and time
- **AND** the time SHALL be formatted correctly
- **AND** support time range display

#### Scenario: Activity location component
- **WHEN** adding activity location
- **THEN** the administrator SHALL be able to enter location details
- **AND** optionally add map integration
- **AND** support address formatting

#### Scenario: Activity highlights component
- **WHEN** adding activity highlights
- **THEN** the administrator SHALL be able to add multiple highlight points
- **AND** each point SHALL support icon selection
- **AND** support reordering of highlights

#### Scenario: Activity details component
- **WHEN** adding activity details
- **THEN** the administrator SHALL be able to enter rich text content
- **AND** support basic text formatting
- **AND** support image insertion

#### Scenario: Ticket information component
- **WHEN** adding ticket information
- **THEN** the administrator SHALL be able to define ticket types
- **AND** set prices for each type
- **AND** specify ticket availability

#### Scenario: Refund policy component
- **WHEN** adding refund policy
- **THEN** the administrator SHALL be able to enter policy text
- **AND** support multiple policy rules
- **AND** format policy display

#### Scenario: Payment and invoice component
- **WHEN** adding payment information
- **THEN** the administrator SHALL be able to configure payment methods
- **AND** specify invoice options
- **AND** add payment instructions

#### Scenario: Organizer information component
- **WHEN** adding organizer information
- **THEN** the administrator SHALL be able to enter organizer details
- **AND** add organizer logo
- **AND** include contact information

#### Scenario: Divider component
- **WHEN** adding a divider
- **THEN** a visual separator SHALL be inserted
- **AND** the administrator SHALL be able to customize divider style
- **AND** adjust divider spacing

#### Scenario: Registration button component
- **WHEN** adding a registration button
- **THEN** the administrator SHALL be able to customize button text
- **AND** set button color and style
- **AND** configure button action

### Requirement: Mobile-Optimized Interaction
The page designer SHALL provide mobile-optimized interaction patterns.

#### Scenario: Component library access
- **WHEN** the administrator wants to add a component
- **THEN** a bottom action sheet SHALL display the component library
- **AND** components SHALL be organized by category
- **AND** each component SHALL have a clear icon and label

#### Scenario: Component editing interface
- **WHEN** editing a component
- **THEN** a bottom popup SHALL display the editor
- **AND** the editor SHALL occupy 70% of screen height
- **AND** support scrolling for long forms
- **AND** have clear save and cancel buttons

#### Scenario: Touch-friendly controls
- **WHEN** interacting with components
- **THEN** all touch targets SHALL be at least 44x44 pixels
- **AND** provide visual feedback on touch
- **AND** support common touch gestures

#### Scenario: Responsive layout
- **WHEN** viewing on different screen sizes
- **THEN** the layout SHALL adapt appropriately
- **AND** maintain usability on small screens
- **AND** optimize for portrait orientation

### Requirement: Data Persistence
The page designer SHALL properly save and load page configurations.

#### Scenario: Configuration data structure
- **WHEN** saving page configuration
- **THEN** it SHALL include all component data
- **AND** preserve component order
- **AND** include theme settings
- **AND** be stored in JSON format

#### Scenario: Validation before save
- **WHEN** saving page configuration
- **THEN** required fields SHALL be validated
- **AND** invalid data SHALL be rejected
- **AND** clear error messages SHALL be displayed

#### Scenario: Auto-save functionality
- **WHEN** the administrator makes changes
- **THEN** changes SHALL be auto-saved periodically
- **AND** save status SHALL be indicated
- **AND** unsaved changes SHALL be warned before exit

### Requirement: Brand admin generates campaign posters
系统 SHALL 允许品牌管理员在管理端为活动生成海报。

#### Scenario: Generate poster from campaign list
- **WHEN** 品牌管理员在活动列表点击“生成海报”
- **THEN** 系统 SHALL 使用该活动的海报模板生成海报
- **AND** 返回可预览/下载的海报 URL

### Requirement: Brand admin configures distribution rules in campaign editor
系统 SHALL 允许品牌管理员在活动创建/编辑时配置分销规则。

#### Scenario: Configure distribution rules
- **WHEN** 品牌管理员启用分销
- **THEN** 系统 SHALL 允许选择分销层级（1-3级）
- **AND** 允许填写各级奖励比例
- **AND** 配置结果 SHALL 保存到活动数据中

