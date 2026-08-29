CREATE TABLE audit_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL CHECK (
        entity_type IN ('client', 'project', 'contract', 'blog_post')
    ),
    entity_id INTEGER,
    summary TEXT NOT NULL,
    details TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_events_created_at
ON audit_events(created_at DESC, id DESC);

CREATE INDEX idx_audit_events_entity
ON audit_events(entity_type, entity_id, created_at DESC);

CREATE INDEX idx_clients_status
ON clients(status);

CREATE INDEX idx_projects_status_due_date
ON projects(status, due_date);

CREATE INDEX idx_contracts_status_end_date
ON contracts(status, end_date);

CREATE INDEX idx_blog_posts_created_at
ON blog_posts(created_at DESC);

-- SQLite cannot add CHECK constraints to existing tables without rebuilding them.
-- These triggers enforce the same invariants for every future direct write while
-- preserving the existing production data and migration history.

CREATE TRIGGER clients_validate_insert
BEFORE INSERT ON clients
WHEN
    trim(NEW.name) = '' OR
    trim(NEW.email) = '' OR
    NEW.status NOT IN ('active', 'inactive')
BEGIN
    SELECT RAISE(ABORT, 'invalid client record');
END;

CREATE TRIGGER clients_validate_update
BEFORE UPDATE ON clients
WHEN
    trim(NEW.name) = '' OR
    trim(NEW.email) = '' OR
    NEW.status NOT IN ('active', 'inactive')
BEGIN
    SELECT RAISE(ABORT, 'invalid client record');
END;

CREATE TRIGGER projects_validate_insert
BEFORE INSERT ON projects
WHEN
    trim(NEW.name) = '' OR
    NEW.status NOT IN ('planned', 'active', 'completed', 'archived') OR
    (NEW.start_date IS NOT NULL AND julianday(NEW.start_date) IS NULL) OR
    (NEW.due_date IS NOT NULL AND julianday(NEW.due_date) IS NULL) OR
    (
        NEW.start_date IS NOT NULL AND
        NEW.due_date IS NOT NULL AND
        julianday(NEW.due_date) < julianday(NEW.start_date)
    )
BEGIN
    SELECT RAISE(ABORT, 'invalid project record');
END;

CREATE TRIGGER projects_validate_update
BEFORE UPDATE ON projects
WHEN
    trim(NEW.name) = '' OR
    NEW.status NOT IN ('planned', 'active', 'completed', 'archived') OR
    (NEW.start_date IS NOT NULL AND julianday(NEW.start_date) IS NULL) OR
    (NEW.due_date IS NOT NULL AND julianday(NEW.due_date) IS NULL) OR
    (
        NEW.start_date IS NOT NULL AND
        NEW.due_date IS NOT NULL AND
        julianday(NEW.due_date) < julianday(NEW.start_date)
    )
BEGIN
    SELECT RAISE(ABORT, 'invalid project record');
END;

CREATE TRIGGER contracts_validate_insert
BEFORE INSERT ON contracts
WHEN
    trim(NEW.title) = '' OR
    NEW.status NOT IN ('draft', 'sent', 'accepted', 'completed', 'cancelled') OR
    (NEW.value_cents IS NOT NULL AND NEW.value_cents < 0) OR
    (NEW.start_date IS NOT NULL AND julianday(NEW.start_date) IS NULL) OR
    (NEW.end_date IS NOT NULL AND julianday(NEW.end_date) IS NULL) OR
    (
        NEW.start_date IS NOT NULL AND
        NEW.end_date IS NOT NULL AND
        julianday(NEW.end_date) < julianday(NEW.start_date)
    ) OR
    (
        NEW.project_id IS NOT NULL AND
        NOT EXISTS (
            SELECT 1
            FROM projects
            WHERE projects.id = NEW.project_id
              AND projects.client_id = NEW.client_id
        )
    )
BEGIN
    SELECT RAISE(ABORT, 'invalid contract record');
END;

CREATE TRIGGER contracts_validate_update
BEFORE UPDATE ON contracts
WHEN
    trim(NEW.title) = '' OR
    NEW.status NOT IN ('draft', 'sent', 'accepted', 'completed', 'cancelled') OR
    (NEW.value_cents IS NOT NULL AND NEW.value_cents < 0) OR
    (NEW.start_date IS NOT NULL AND julianday(NEW.start_date) IS NULL) OR
    (NEW.end_date IS NOT NULL AND julianday(NEW.end_date) IS NULL) OR
    (
        NEW.start_date IS NOT NULL AND
        NEW.end_date IS NOT NULL AND
        julianday(NEW.end_date) < julianday(NEW.start_date)
    ) OR
    (
        NEW.project_id IS NOT NULL AND
        NOT EXISTS (
            SELECT 1
            FROM projects
            WHERE projects.id = NEW.project_id
              AND projects.client_id = NEW.client_id
        )
    )
BEGIN
    SELECT RAISE(ABORT, 'invalid contract record');
END;

CREATE TRIGGER blog_posts_validate_insert
BEFORE INSERT ON blog_posts
WHEN
    trim(NEW.title) = '' OR
    trim(NEW.slug) = '' OR
    trim(NEW.excerpt) = '' OR
    trim(NEW.content) = '' OR
    NEW.status NOT IN ('draft', 'published') OR
    (NEW.status = 'draft' AND NEW.published_at IS NOT NULL) OR
    (NEW.status = 'published' AND NEW.published_at IS NULL)
BEGIN
    SELECT RAISE(ABORT, 'invalid blog post record');
END;

CREATE TRIGGER blog_posts_validate_update
BEFORE UPDATE ON blog_posts
WHEN
    trim(NEW.title) = '' OR
    trim(NEW.slug) = '' OR
    trim(NEW.excerpt) = '' OR
    trim(NEW.content) = '' OR
    NEW.status NOT IN ('draft', 'published') OR
    (NEW.status = 'draft' AND NEW.published_at IS NOT NULL) OR
    (NEW.status = 'published' AND NEW.published_at IS NULL)
BEGIN
    SELECT RAISE(ABORT, 'invalid blog post record');
END;

CREATE TRIGGER clients_audit_insert
AFTER INSERT ON clients
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary
    ) VALUES (
        'client.created',
        'client',
        NEW.id,
        'Created client: ' || NEW.name
    );
END;

CREATE TRIGGER clients_audit_update
AFTER UPDATE ON clients
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary,
        details
    ) VALUES (
        CASE
            WHEN OLD.status IS NOT NEW.status THEN 'client.status_changed'
            ELSE 'client.updated'
        END,
        'client',
        NEW.id,
        CASE
            WHEN OLD.status IS NOT NEW.status THEN 'Changed client status: ' || NEW.name
            ELSE 'Updated client: ' || NEW.name
        END,
        CASE
            WHEN OLD.status IS NOT NEW.status THEN OLD.status || ' → ' || NEW.status
            ELSE ''
        END
    );
END;

CREATE TRIGGER clients_audit_delete
AFTER DELETE ON clients
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary
    ) VALUES (
        'client.deleted',
        'client',
        OLD.id,
        'Deleted client: ' || OLD.name
    );
END;

CREATE TRIGGER projects_audit_insert
AFTER INSERT ON projects
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary
    ) VALUES (
        'project.created',
        'project',
        NEW.id,
        'Created project: ' || NEW.name
    );
END;

CREATE TRIGGER projects_audit_update
AFTER UPDATE ON projects
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary,
        details
    ) VALUES (
        CASE
            WHEN OLD.status IS NOT NEW.status AND NEW.status = 'archived' THEN 'project.archived'
            WHEN OLD.status IS NOT NEW.status THEN 'project.status_changed'
            ELSE 'project.updated'
        END,
        'project',
        NEW.id,
        CASE
            WHEN OLD.status IS NOT NEW.status AND NEW.status = 'archived' THEN 'Archived project: ' || NEW.name
            WHEN OLD.status IS NOT NEW.status THEN 'Changed project status: ' || NEW.name
            ELSE 'Updated project: ' || NEW.name
        END,
        CASE
            WHEN OLD.status IS NOT NEW.status THEN OLD.status || ' → ' || NEW.status
            ELSE ''
        END
    );
END;

CREATE TRIGGER projects_audit_delete
AFTER DELETE ON projects
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary
    ) VALUES (
        'project.deleted',
        'project',
        OLD.id,
        'Deleted project: ' || OLD.name
    );
END;

CREATE TRIGGER contracts_audit_insert
AFTER INSERT ON contracts
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary
    ) VALUES (
        'contract.created',
        'contract',
        NEW.id,
        'Created contract: ' || NEW.title
    );
END;

CREATE TRIGGER contracts_audit_update
AFTER UPDATE ON contracts
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary,
        details
    ) VALUES (
        CASE
            WHEN OLD.status IS NOT NEW.status AND NEW.status = 'cancelled' THEN 'contract.cancelled'
            WHEN OLD.status IS NOT NEW.status THEN 'contract.status_changed'
            ELSE 'contract.updated'
        END,
        'contract',
        NEW.id,
        CASE
            WHEN OLD.status IS NOT NEW.status AND NEW.status = 'cancelled' THEN 'Cancelled contract: ' || NEW.title
            WHEN OLD.status IS NOT NEW.status THEN 'Changed contract status: ' || NEW.title
            ELSE 'Updated contract: ' || NEW.title
        END,
        CASE
            WHEN OLD.status IS NOT NEW.status THEN OLD.status || ' → ' || NEW.status
            ELSE ''
        END
    );
END;

CREATE TRIGGER contracts_audit_delete
AFTER DELETE ON contracts
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary
    ) VALUES (
        'contract.deleted',
        'contract',
        OLD.id,
        'Deleted contract: ' || OLD.title
    );
END;

CREATE TRIGGER blog_posts_audit_insert
AFTER INSERT ON blog_posts
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary
    ) VALUES (
        CASE
            WHEN NEW.status = 'published' THEN 'blog_post.published'
            ELSE 'blog_post.created'
        END,
        'blog_post',
        NEW.id,
        CASE
            WHEN NEW.status = 'published' THEN 'Published journal post: ' || NEW.title
            ELSE 'Created journal draft: ' || NEW.title
        END
    );
END;

CREATE TRIGGER blog_posts_audit_update
AFTER UPDATE ON blog_posts
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary,
        details
    ) VALUES (
        CASE
            WHEN OLD.status = 'draft' AND NEW.status = 'published' THEN 'blog_post.published'
            WHEN OLD.status = 'published' AND NEW.status = 'draft' THEN 'blog_post.unpublished'
            ELSE 'blog_post.updated'
        END,
        'blog_post',
        NEW.id,
        CASE
            WHEN OLD.status = 'draft' AND NEW.status = 'published' THEN 'Published journal post: ' || NEW.title
            WHEN OLD.status = 'published' AND NEW.status = 'draft' THEN 'Unpublished journal post: ' || NEW.title
            ELSE 'Updated journal post: ' || NEW.title
        END,
        CASE
            WHEN OLD.status IS NOT NEW.status THEN OLD.status || ' → ' || NEW.status
            ELSE ''
        END
    );
END;

CREATE TRIGGER blog_posts_audit_delete
AFTER DELETE ON blog_posts
BEGIN
    INSERT INTO audit_events (
        action,
        entity_type,
        entity_id,
        summary
    ) VALUES (
        'blog_post.deleted',
        'blog_post',
        OLD.id,
        'Deleted journal post: ' || OLD.title
    );
END;
