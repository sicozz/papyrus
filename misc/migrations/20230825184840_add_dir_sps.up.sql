CREATE OR REPLACE FUNCTION sp_getDepth(dirUuid UUID) RETURNS INT AS $PROC$
DECLARE
    rootUuid UUID;
    u UUID;
    depth INT := 0;
BEGIN
    rootUuid := '00000000-0000-0000-0000-000000000000';
    u := dirUuid;

    WHILE u <> rootUuid LOOP
        SELECT parent_dir INTO u FROM dir WHERE uuid = u;
        depth := depth + 1;
    END LOOP;

    RETURN depth;
END;
$PROC$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sp_getPath(dirUuid UUID) RETURNS TEXT AS $PROC$
DECLARE
    rootUuid UUID;
    u UUID;
    dirName TEXT;
    path TEXT;
BEGIN
    rootUuid := '00000000-0000-0000-0000-000000000000';
    u := dirUuid;
    path := '';

    IF rootUuid = u THEN
        RETURN '/';
    END IF;

    WHILE u <> rootUuid LOOP
        SELECT parent_dir, name INTO u, dirName FROM dir WHERE uuid = u;
        path := '/' || dirName || path;
    END LOOP;

    RETURN path;
END;
$PROC$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sp_isSubDir(dirUuid UUID, destDir UUID) RETURNS BOOLEAN AS $PROC$
DECLARE
    rootUuid UUID;
    u UUID;
    v UUID;
    res BOOLEAN := FALSE;
BEGIN
    rootUuid := '00000000-0000-0000-0000-000000000000';
    u := dirUuid;

    WHILE u <> rootUuid LOOP
        SELECT parent_dir INTO v FROM dir WHERE uuid = u;
        IF destDir = v THEN
            res := TRUE;
            u := rootUuid;
        ELSE
            u := v;
        END IF;
    END LOOP;

    RETURN res;
END;
$PROC$ LANGUAGE plpgsql;
