import clsx from "clsx";
import {User} from "../context/Chat.context";
import {supportedLanguages} from "./LanguageSelector";

interface UsersListProps {
    users: User[];
    className?: string;
}
const UsersList: React.FC<UsersListProps> = ({users, className}) => {
    return (
        <div className={clsx("flex w-full flex-col bg-gray-50 text-theme-base gap-2", className)}>
            {users.map(({id, username, language}) => {
                return (
                    <div key={id} className="px-4 py-2 flex gap-2 rounded">
                        <span>{supportedLanguages[language].icon}</span>
                        <span>{username}</span>
                    </div>
                );
            })}
        </div>
    );
};
export default UsersList;
