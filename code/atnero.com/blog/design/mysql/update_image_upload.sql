USE blog;
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'upload_image_100K',
        '上传100KB的图片',
        true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'upload_image_400K',
        '上传400KB的图片',
        true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'upload_image_1M',
        '上传1MB的图片',
        false);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'upload_image_10M',
        '上传10MB的图片',
        false);
INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'base_user_rightset' AND user_right_item.name = 'upload_image_100K';
INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'base_user_rightset' AND user_right_item.name = 'upload_image_400K';