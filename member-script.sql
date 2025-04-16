USE [member_service]
GO
/****** Object:  Table [dbo].[sys_channel]    Script Date: 16/4/2025 11:15:21 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[sys_channel](
	[id] [bigint] IDENTITY(1,1) NOT NULL,
	[app_id] [varchar](100) NOT NULL,
	[app_key] [varchar](100) NOT NULL,
	[status] [char](2) NOT NULL,
	[sig_method] [varchar](100) NOT NULL,
	[create_time] [datetime] NULL,
	[update_time] [datetime] NOT NULL,
PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[users]    Script Date: 16/4/2025 11:15:22 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[users](
	[id] [bigint] IDENTITY(18,1) NOT NULL,
	[external_id] [varchar](255) NULL,
	[external_id_type] [varchar](50) NULL,
	[email] [varchar](50) NULL,
	[burn_pin] [int] NULL,
	[session_token] [varchar](255) NULL,
	[session_expiry] [int] NULL,
	[gr_id] [varchar](255) NULL,
	[rlp_id] [varchar](255) NULL,
	[rws_membership_id] [varchar](255) NULL,
	[rws_membership_number] [int] NULL,
	[created_at] [datetime] NOT NULL,
	[updated_at] [datetime] NOT NULL,
 CONSTRAINT [PK_users] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[sys_channel] ADD  DEFAULT ('10') FOR [status]
GO
ALTER TABLE [dbo].[sys_channel] ADD  DEFAULT ('SHA256') FOR [sig_method]
GO
ALTER TABLE [dbo].[users] ADD  DEFAULT (getdate()) FOR [created_at]
GO
ALTER TABLE [dbo].[users] ADD  DEFAULT (getdate()) FOR [updated_at]
GO
