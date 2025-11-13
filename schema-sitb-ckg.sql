-- SITB CKG Database Schema
-- This schema supports the CKG (Cegah Komplikasi Gizi) screening system

-- Enable strict SQL mode for better data integrity
SET SQL_MODE = "STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION";

-- Set default character set and collation
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =============================================
-- DATABASE CREATION
-- =============================================
-- Create database if not exists
-- CREATE DATABASE IF NOT EXISTS `sitb_ckg` 
-- CHARACTER SET utf8mb4 
-- COLLATE utf8mb4_unicode_ci;

-- USE `sitb_ckg`;

-- =============================================
-- TABLE: tmp_ckg_incoming (Pub/Sub incoming messages - untuk SITB)
-- =============================================
CREATE TABLE `tmp_ckg_incoming` (
    `id` VARCHAR(100) NOT NULL COMMENT 'Message ID from Pub/Sub',
    `data` JSON NOT NULL COMMENT 'Message data in JSON format',
    `attributes` JSON NULL COMMENT 'Message attributes in JSON format',
    `received_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Message received timestamp',
    PRIMARY KEY (`id`),
    INDEX `idx_received_at` (`received_at`)
    -- INDEX `idx_processed` (`processed`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Pub/Sub Incoming Messages Table';

-- =============================================
-- TABLE: ckg_pubsub_incoming (Pub/Sub incoming messages - untuk CKG)
-- =============================================
CREATE TABLE `ckg_pubsub_incoming` (
    `id` VARCHAR(100) NOT NULL COMMENT 'Message ID from Pub/Sub',
    `data` JSON NOT NULL COMMENT 'Message data in JSON format',
    `received_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Message received timestamp',
    `processed_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Message received timestamp',
    PRIMARY KEY (`id`),
    INDEX `idx_received_at` (`received_at`),
    INDEX `idx_processed_at` (`processed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Pub/Sub Incoming Messages Table';

-- =============================================
-- TABLE: tmp_ckg_outgoing (API outgoing messages - untuk SITB)
-- =============================================
CREATE TABLE `tmp_ckg_outgoing` (
    `terduga_id` VARCHAR(50) NOT NULL COMMENT 'Suspected case ID',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update timestamp',
    PRIMARY KEY (`terduga_id`),
    INDEX `idx_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API Outgoing Messages Table';

-- =============================================
-- TABLE: ckg_pubsu_outgoing (API outgoing messages - untuk CKG)
-- =============================================
CREATE TABLE `ckg_pubsu_outgoing` (
    `id` VARCHAR(100) NOT NULL COMMENT 'Message ID from Pub/Sub',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record create timestamp',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update timestamp',
    PRIMARY KEY (`id`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API Outgoing Messages Table';

-- =============================================
-- END OF SCHEMA
-- =============================================
SET FOREIGN_KEY_CHECKS = 1;

-- Schema creation completed
-- Database: sitb_ckg
-- Tables: 3 (tmp_ckg_incoming, tmp_ckg_processed, tmp_ckg_outgoing)