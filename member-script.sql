USE [master]
GO
/****** Object:  Database [member_service]    Script Date: 16/4/2025 11:19:51 PM ******/
CREATE DATABASE [member_service]
 CONTAINMENT = NONE
 ON  PRIMARY 
( NAME = N'member_service', FILENAME = N'/var/opt/mssql/data/member_service.mdf' , SIZE = 8192KB , MAXSIZE = UNLIMITED, FILEGROWTH = 65536KB )
 LOG ON 
( NAME = N'member_service_log', FILENAME = N'/var/opt/mssql/data/member_service_log.ldf' , SIZE = 8192KB , MAXSIZE = 2048GB , FILEGROWTH = 65536KB )
 WITH CATALOG_COLLATION = DATABASE_DEFAULT
GO
ALTER DATABASE [member_service] SET COMPATIBILITY_LEVEL = 150
GO
IF (1 = FULLTEXTSERVICEPROPERTY('IsFullTextInstalled'))
begin
EXEC [member_service].[dbo].[sp_fulltext_database] @action = 'enable'
end
GO
ALTER DATABASE [member_service] SET ANSI_NULL_DEFAULT OFF 
GO
ALTER DATABASE [member_service] SET ANSI_NULLS OFF 
GO
ALTER DATABASE [member_service] SET ANSI_PADDING OFF 
GO
ALTER DATABASE [member_service] SET ANSI_WARNINGS OFF 
GO
ALTER DATABASE [member_service] SET ARITHABORT OFF 
GO
ALTER DATABASE [member_service] SET AUTO_CLOSE OFF 
GO
ALTER DATABASE [member_service] SET AUTO_SHRINK OFF 
GO
ALTER DATABASE [member_service] SET AUTO_UPDATE_STATISTICS ON 
GO
ALTER DATABASE [member_service] SET CURSOR_CLOSE_ON_COMMIT OFF 
GO
ALTER DATABASE [member_service] SET CURSOR_DEFAULT  GLOBAL 
GO
ALTER DATABASE [member_service] SET CONCAT_NULL_YIELDS_NULL OFF 
GO
ALTER DATABASE [member_service] SET NUMERIC_ROUNDABORT OFF 
GO
ALTER DATABASE [member_service] SET QUOTED_IDENTIFIER OFF 
GO
ALTER DATABASE [member_service] SET RECURSIVE_TRIGGERS OFF 
GO
ALTER DATABASE [member_service] SET  DISABLE_BROKER 
GO
ALTER DATABASE [member_service] SET AUTO_UPDATE_STATISTICS_ASYNC OFF 
GO
ALTER DATABASE [member_service] SET DATE_CORRELATION_OPTIMIZATION OFF 
GO
ALTER DATABASE [member_service] SET TRUSTWORTHY OFF 
GO
ALTER DATABASE [member_service] SET ALLOW_SNAPSHOT_ISOLATION OFF 
GO
ALTER DATABASE [member_service] SET PARAMETERIZATION SIMPLE 
GO
ALTER DATABASE [member_service] SET READ_COMMITTED_SNAPSHOT OFF 
GO
ALTER DATABASE [member_service] SET HONOR_BROKER_PRIORITY OFF 
GO
ALTER DATABASE [member_service] SET RECOVERY FULL 
GO
ALTER DATABASE [member_service] SET  MULTI_USER 
GO
ALTER DATABASE [member_service] SET PAGE_VERIFY CHECKSUM  
GO
ALTER DATABASE [member_service] SET DB_CHAINING OFF 
GO
ALTER DATABASE [member_service] SET FILESTREAM( NON_TRANSACTED_ACCESS = OFF ) 
GO
ALTER DATABASE [member_service] SET TARGET_RECOVERY_TIME = 60 SECONDS 
GO
ALTER DATABASE [member_service] SET DELAYED_DURABILITY = DISABLED 
GO
ALTER DATABASE [member_service] SET ACCELERATED_DATABASE_RECOVERY = OFF  
GO
EXEC sys.sp_db_vardecimal_storage_format N'member_service', N'ON'
GO
ALTER DATABASE [member_service] SET QUERY_STORE = OFF
GO
USE [member_service]
GO
/****** Object:  Table [dbo].[sys_channel]    Script Date: 16/4/2025 11:19:51 PM ******/
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
/****** Object:  Table [dbo].[users]    Script Date: 16/4/2025 11:19:52 PM ******/
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
USE [master]
GO
ALTER DATABASE [member_service] SET  READ_WRITE 
GO
