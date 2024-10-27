/*
 * Copyright (C) 2024 Mustafa Naseer
 *
 * This file is part of urulink chat application.
 *
 * urulink is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation version 3 of the License .
 *
 *
 * urulink is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with urulink. If not, see <http://www.gnu.org/licenses/>.
 */

package models

type ClientsLoginResponse struct {
	Uid      string `json:"uid"`
	Username string `json:"username"`
}

type FileSender struct {
	FileName string `json:"file_name"`
	FileUrl  string `json:"file_url"`
}

type FileTypeConfig struct {
	Dir     string
	MaxSize int64
}