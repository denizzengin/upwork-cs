<!--
 Copyright 2022 TCDZENGIN

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
-->

# upwork-cs

Upwork case study.

## Installation

Clone the github repository.

```npm

INSERT INTO public.roles (created_at, updated_at, is_deleted, role_name) VALUES(now(),now(), false, 'admin');
INSERT INTO public.roles (created_at, updated_at, is_deleted, role_name) VALUES(now(),now(), false, 'user');

INSERT INTO public.role_rights (created_at, updated_at, is_deleted, pattern, role_id) VALUES(now(),now(), false, '/api/users/me', 2);
INSERT INTO public.role_rights (created_at, updated_at, is_deleted, pattern, role_id) VALUES(now(),now(), false, '/api/users/me', 1);
INSERT INTO public.role_rights (created_at, updated_at, is_deleted, pattern, role_id) VALUES(now(),now(), false, '/api/users', 1);
INSERT INTO public.role_rights (created_at, updated_at, is_deleted, pattern, role_id) VALUES(now(),now(), false, '/api/users/admin', 1);
INSERT INTO public.role_rights (created_at, updated_at, is_deleted, pattern, role_id) VALUES(now(),now(), false, '/api/users/me', 1);

$ git clone https://github.com/denizzengin/upwork-cs.git
$ make build
$ make run
$ docker-compose up --build
```
